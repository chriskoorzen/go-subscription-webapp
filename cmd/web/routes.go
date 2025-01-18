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

	mux.Mount("/members", app.authRouter())

	return mux
}

func (app *Config) authRouter() http.Handler {
	// create a new chi router
	mux := chi.NewRouter()

	// set up middleware
	mux.Use(app.Auth)

	// set up protected routes
	mux.Get("/plans", app.GETSubscriptionPlans)
	// mux.Post("/subscribe", app.POSTSubscribePage)

	return mux
}
