package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	// create a new chi router
	mux := chi.NewRouter()

	// set up middleware
	mux.Use(middleware.Recoverer) // recover from panics
	mux.Use(app.SessionLoad)      // load and save session data

	// set up routes
	mux.Get("/", app.GETHomePage)
	mux.Get("/login", app.GETLoginPage)
	mux.Post("/login", app.POSTLoginPage)
	mux.Get("/logout", app.GETLogout)
	mux.Get("/register", app.GETRegisterPage)
	mux.Post("/register", app.POSTRegisterPage)
	mux.Get("/activate-account", app.GETActivateAccount)

	// testing email route
	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Println("[ DEV ] GET /test-email")
		app.InfoLog.Println("[ DEV ] Firing off email send function")
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromName:    "Test",
			FromAddress: "test@mycompany.com",
			ErrorChan:   make(chan error),
		}

		msg := Message{
			To:      "me@testdomain.com",
			Subject: "Test email",
			Data:    "Hello! This is a test email.",
		}

		m.Send(msg, make(chan error))
	})

	return mux
}
