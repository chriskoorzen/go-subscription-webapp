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

	return mux
}
