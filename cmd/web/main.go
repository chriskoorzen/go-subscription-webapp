package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	fmt.Println("Hello, Subscription Service!")

	// connect to database
	database := initDB()

	// create sessions
	session := initSession()

	// Create loggers
	infoLog := log.New(os.Stdout, color.GreenString("[INFO\t] "), log.Ldate|log.Ltime)
	successLog := log.New(os.Stdout, color.CyanString("[SUCCESS] "), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, color.RedString("[ERROR\t] "), log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// Create waitgroup - part of the graceful shutdown protocol.
	// Spawned background processes increment the waitgroup counter, and decrement it when they finish.
	// When the server receives a shutdown signal, it waits for the waitgroup counter to reach 0
	wg := sync.WaitGroup{}

	// setup app config
	app := Config{
		Session:    session,
		DB:         database,
		Wait:       &wg,
		InfoLog:    infoLog,
		SuccessLog: successLog,
		ErrorLog:   errorLog,
		Models:     db.New(database),
	}

	// set up mail
	app.Mailer = app.initMailer()
	go app.listenForMail()

	// listen for shutdown signal
	go app.listenForShutDown()

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
	// register custom types
	gob.Register(db.User{})

	// Create a new session manager and store it in the Config struct
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

func (app *Config) initMailer() Mail {

	m := Mail{
		// TODO - get these from environment variables
		Domain:     "localhost",
		Host:       "localhost",
		Port:       1025,
		Encryption: "none",
		// Username:    "",
		// Password:    "",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",

		Wait:       app.Wait,                // same as the app waitgroup (global waitgroup)
		MailerChan: make(chan Message, 100), // buffered channel
		ErrorChan:  make(chan error),
		DoneChan:   make(chan bool),
	}

	return m
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

func (app *Config) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit                                      // block until signal is received
	app.InfoLog.Println("Starting shutdown...") // log message
	app.shutdown()                              // call shutdown function
	os.Exit(0)                                  // exit gracefully
}

func (app *Config) shutdown() {
	// Do cleanup here
	app.InfoLog.Println("Waiting for background processes to finish...")

	// block until waitgroup counter is 0
	app.Wait.Wait()

	app.InfoLog.Println("All background processes finished. Shutting down...")
}
