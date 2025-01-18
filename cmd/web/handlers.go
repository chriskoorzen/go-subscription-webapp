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

	app.Session.RenewToken(r.Context()) // renew the session token when logging in

	// parse POST form
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println("Error parsing form: ", err)
		http.Error(w, "Something went wrong. Please try again.", http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// authenticate user
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials") // store error message in session
		app.ErrorLog.Println("Error getting user by email: ", err)   // log error
		http.Redirect(w, r, "/login", http.StatusSeeOther)           // redirect back to login page
		return
	}

	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials") // store error message in session
		app.ErrorLog.Println("Error comparing passwords: ", err)     // log error
		http.Redirect(w, r, "/login", http.StatusSeeOther)           // redirect back to login page
		return
	}
	if !validPassword {
		app.Session.Put(r.Context(), "error", "Invalid credentials") // store error message in session
		app.ErrorLog.Println("Invalid password")                     // log error
		http.Redirect(w, r, "/login", http.StatusSeeOther)           // redirect back to login page
		return
	}

	// Auth passed - Log in user
	app.Session.Put(r.Context(), "userID", user.ID) // store user ID in session
	app.Session.Put(r.Context(), "user", user)      // store user data in session
	app.Session.Put(r.Context(), "flash", "You've been logged in successfully")

	app.SuccessLog.Printf("User %d logged in", user.ID)
	// Redirect to "successs page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) GETLogout(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)

	userID := app.Session.GetInt(r.Context(), "userID")

	// Clean up session
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	app.SuccessLog.Printf("User %d logged out", userID)
	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
