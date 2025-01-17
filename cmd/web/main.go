package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	fmt.Println("Hello, Subscription Service!")

	// connect to database
	db := initDB()
	fmt.Println(db)

	// create sessions

	// create channels

	// create waitgroup

	// setup app config

	// set up mail

	// listen for connections
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
