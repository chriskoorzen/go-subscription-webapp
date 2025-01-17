package main

import (
	"net/http"
)

func (app *Config) GETHomePage(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)
	app.render(w, r, "home.page.gohtml", nil)
}
