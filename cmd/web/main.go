package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	fmt.Println("Hello, Subscription Service!")

	// connect to database
	db := initDB()

	// create sessions
	session := initSession()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// setup app config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	// set up mail

	// listen for connections
	app.serve()
}

func initDB() *sql.DB {

	conn := connectToDB()
	if conn == nil {
		log.Panic("Failed to connect to database")
	}

	return conn
}

func connectToDB() *sql.DB {
	counts := 0 // number of attempts to connect to database

	dsn := os.Getenv("DSN") // Database Source Name

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Database not yet ready...")
		} else {
			log.Println("Connected to database")
			return connection
		}

		// Retry every 1 second
		if counts > 5 {
			return nil
		}
		counts++
		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// Open connection to database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Check if connection is valid
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {
	session := scs.New()
	session.Store = redisstore.New(initRedis()) // Use Redis to store session data

	// Set session options
	session.Lifetime = 24 * time.Hour              // 24 hours before session expires
	session.Cookie.Persist = true                  // Persist session even after browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode // SameSite cookie policy
	session.Cookie.Secure = true                   // Secure cookie policy

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS")) // Redis connection string
		},
	}
	return redisPool
}

func (app *Config) serve() {
	// start http server
	server := &http.Server{
		Addr:    ":8811",
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting server on port :8811")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
