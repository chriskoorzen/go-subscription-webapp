package main

import (
	"net/http"
)

func (app *Config) GETHomePage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) GETLoginPage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) POSTLoginPage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("POST %s\n", r.URL.Path)
	// TODO
}

func (app *Config) GETLogout(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	// TODO
}

func (app *Config) GETRegisterPage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) POSTRegisterPage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("POST %s\n", r.URL.Path)
	// TODO
}

// Sent once the user has successfully registered
// so we can verify their email address
func (app *Config) GETActivateAccount(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	// TODO
}
