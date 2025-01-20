package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
	"github.com/fatih/color"
)

var testApp Config

func TestMain(m *testing.M) {
	// register custom types for session
	gob.Register(db.User{})

	// override paths for testing
	pathToManual = "../../pdf"
	pathToTmpPDFWrite = "../../tmp"
	pathToTemplates = "./templates"

	// Set up session
	session := scs.New()
	session.Lifetime = 24 * time.Hour              // 24 hours before session expires
	session.Cookie.Persist = true                  // Persist session even after browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode // SameSite cookie policy
	session.Cookie.Secure = true                   // Secure cookie policy

	// Set up testApp environment
	testApp = Config{
		Session:       session,
		DB:            nil,             // do not connect to database for this test
		Models:        db.TestNew(nil), // "database free" models
		InfoLog:       log.New(os.Stdout, color.GreenString("[INFO\t] "), log.Ldate|log.Ltime),
		SuccessLog:    log.New(os.Stdout, color.CyanString("[SUCCESS] "), log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, color.RedString("[ERROR\t] "), log.Ldate|log.Ltime|log.Lshortfile),
		Wait:          &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	// create a dummy mailer
	testApp.Mailer = Mail{
		ErrorChan:  make(chan error),
		MailerChan: make(chan Message, 100),
		DoneChan:   make(chan bool),
		Wait:       testApp.Wait,
	}

	go func() { // simply consume the mailer channels
		for {
			select {
			case <-testApp.Mailer.MailerChan:
				testApp.Wait.Done() // decrement the waitgroup for sending mail
			case <-testApp.Mailer.ErrorChan:
			case <-testApp.Mailer.DoneChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErrorChan: // output any errors
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return
			}
		}
	}()

	// Run the tests
	os.Exit(m.Run())
}

func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("Error loading session", err)
	}
	return ctx
}
