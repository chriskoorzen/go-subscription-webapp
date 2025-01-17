package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	CSRFToken     string
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	// hardcode template partials in render function
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	var templateFiles = make([]string, 0, len(partials)+1)
	templateFiles = append(templateFiles, fmt.Sprintf("%s/%s", pathToTemplates, t)) // append the template to the slice
	templateFiles = append(templateFiles, partials...)                              // append the partials to the slice

	if td == nil {
		td = &TemplateData{} // if no data is passed in, create an empty TemplateData struct
	}

	// parse the template files
	templ, err := template.ParseFiles(templateFiles...)
	if err != nil {
		app.ErrorLog.Println("Error parsing template files: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the template
	if err := templ.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println("Error executing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")     // get the flash message from the session
	td.Warning = app.Session.PopString(r.Context(), "warning") // get the warning message from the session
	td.Error = app.Session.PopString(r.Context(), "error")     // get the error message from the session
	if app.IsAuthenticated(r) {
		td.Authenticated = true
		// TODO Get other user info and add it to the template data
	}

	return td
}

func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID") // check if the userID key exists in the session
}
