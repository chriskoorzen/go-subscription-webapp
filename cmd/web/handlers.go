package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
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

	if user.Active == 0 {
		app.ErrorLog.Printf("User account %d not activated\n", user.ID)
		app.Session.Put(r.Context(), "error", "Account not activated") // store error message in session
		http.Redirect(w, r, "/login", http.StatusSeeOther)             // redirect back to login page
		return
	}

	validPassword, err := app.Models.User.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials") // store error message in session
		app.ErrorLog.Println("Error comparing passwords: ", err)     // log error
		http.Redirect(w, r, "/login", http.StatusSeeOther)           // redirect back to login page
		return
	}
	if !validPassword {
		// send user a notification email that their account was accessed
		msg := Message{
			To:      email,
			Subject: "Failed login attempt",
			Data:    "Someone tried to log into your account with an incorrect password.",
		}
		app.sendEmail(msg)

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

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println("Error parsing form: ", err)
		app.Session.Put(r.Context(), "error", "Unable to create account.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// TODO Validate form data

	// create a new user
	u := db.User{
		Email:     r.PostForm.Get("email"),
		FirstName: r.PostForm.Get("first-name"),
		LastName:  r.PostForm.Get("last-name"),
		Password:  r.PostForm.Get("password"),
		IsAdmin:   0,
		Active:    0,
	}

	userID, err := app.Models.User.Insert(u)
	if err != nil {
		app.ErrorLog.Println("Error inserting user: ", err)
		app.Session.Put(r.Context(), "error", "Unable to create account.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send activation email
	url := fmt.Sprintf("%s/activate-account?email=%s", "http://localhost:8811", u.Email) // TODO - get this from environment variable
	signedURL := GenerateTokenFromString(url)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTMLEscapeString(signedURL),
	}
	app.sendEmail(msg)

	app.Session.Put(r.Context(), "flash", "Account created. Please check your email to activate your account.")
	app.SuccessLog.Println("User created with ID: ", userID)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Sent once the user has successfully registered
// so we can verify their email address
func (app *Config) GETActivateAccount(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)

	// validate url
	url := r.RequestURI
	testUrl := fmt.Sprintf("%s%s", "http://localhost:8811", url) // TODO - get this from environment variable
	okay := VerifyToken(testUrl)

	if !okay {
		app.ErrorLog.Println("Invalid activation token")
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// activate account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.ErrorLog.Println("Error getting user by email: ", err)
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	u.Active = 1
	err = app.Models.User.Update(*u)
	if err != nil {
		app.ErrorLog.Println("Unable to update user ", err)
		app.Session.Put(r.Context(), "error", "Activation failed")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// success
	app.SuccessLog.Printf("User %d activated account", u.ID)
	app.Session.Put(r.Context(), "flash", "Account activated. Please log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Protected route
func (app *Config) GETSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)

	// get plans
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println("Error getting plans: ", err)
		app.Session.Put(r.Context(), "error", "Unable to get plans")
		// TODO implement error page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

// Protected route
func (app *Config) GETSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Printf("GET %s\n", r.URL.Path)

	// get id of chosen plan
	id := r.URL.Query().Get("plan")
	planID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorLog.Println("Error getting plan: ", err)
		app.Session.Put(r.Context(), "error", "Unable to get plan")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.ErrorLog.Println("Error getting plan: ", err)
		app.Session.Put(r.Context(), "error", "Unable to get plan")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	// get user from session
	user, ok := app.Session.Get(r.Context(), "user").(db.User)
	if !ok {
		app.ErrorLog.Println("Error getting user from session")
		app.Session.Put(r.Context(), "error", "Log in to access this page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice and send email with invoice attached
	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		invoice, err := app.GenerateInvoice(user, plan)
		if err != nil {
			app.ErrorChan <- fmt.Errorf("error generating invoice: %v", err)
		}

		// send email with invoice attached
		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice",
			Template: "invoice-email",
			Data:     invoice,
		}
		app.sendEmail(msg)
	}()

	// generate a manual and send email with manual attached
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()

		pdf := app.GenerateManual(user, plan)
		filePath := fmt.Sprintf("%s/%s_%d_manual.pdf", pathToTmpPDFWrite, time.Now().Format("2006-01-02-15:04"), user.ID)
		err := pdf.OutputFileAndClose(filePath)
		if err != nil {
			app.ErrorChan <- fmt.Errorf("error generating manual: %v", err)
			return
		}
		msg := Message{
			To:      user.Email,
			Subject: "Your Manual",
			Data:    "Please find your manual attached.",
			AttachmentMap: map[string]string{
				"manual.pdf": filePath,
			},
		}
		app.sendEmail(msg)
	}()

	// subscribe user to plan
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.ErrorLog.Println("Error subscribing user to plan: ", err)
		app.Session.Put(r.Context(), "error", "Unable to subscribe to plan")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
	// update user in session
	u, err := app.Models.User.GetOne(user.ID) // get fresh data from db
	if err != nil {
		app.ErrorLog.Println("Error getting user: ", err)
		app.Session.Put(r.Context(), "error", "Unable to get user")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "user", u) // update user in session

	// redirect to success page
	app.Session.Put(r.Context(), "flash", "Subscribed successfully")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}
