package handlers

import (
	"github.com/gorilla/mux"

	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/web/utils"
)

// RegisterAccountHandlers registers all the handlers for the
// endpoints beginning with /account/
func RegisterAccountHandlers(router *mux.Router) {
	router.Handle("/account", utils.Handler(handleGetAccount)).Methods("GET")
	router.Handle("/account", utils.Handler(handlePostAccount)).Methods("POST")
}

// handleGetAccount is called for GET requests at /account.
// If the query string for "action" is "signup", it renders the signup form.
// If the query string for "action" is "signin", it renders the signin form.
// In all the other cases, it renders an error.
func handleGetAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	switch c.R.URL.Query().Get("action") {
	case "signup":
		utils.RenderPage(c, "account/signup", nil)
	case "signin":
		utils.RenderPage(c, "account/signin", nil)
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}

// handlePostAccount is called for POST requests at /account
// If the request value for "action" is "signup", it tries to sign up the user.
// If the request value for "action" is "signin", it tries to sign in the user.
// In all the other cases, it returns an error.
func handlePostAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the request body
	switch c.R.FormValue("action") {
	case "signup":
		signUpUser(c)
	case "signin":
		signInUser(c)
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}

// signUpUser tries to sign up the user whose data is in the request
func signUpUser(c utils.Context) {
	// Fetches data
	user := &database.User{
		Username: c.R.FormValue("username"),
		Email:    c.R.FormValue("email"),
		Password: c.R.FormValue("password"),
	}

	// Ensures the two passwords are equal
	if user.Password != c.R.FormValue("password2") {
		utils.ShowError(c, "Le due password non corrispondono")
		return
	}

	// Tries to sign up the user
	if err := user.SignUp(); err != nil {
		utils.ShowError(c, err.Error())
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// signInUser tries to sign in the user whose data is in the request
func signInUser(c utils.Context) {
	// Fetches data
	user := &database.User{
		Username: c.R.FormValue("username"),
		Password: c.R.FormValue("password"),
	}

	// Tries to sign up the user
	if err := user.SignIn(); err != nil {
		utils.ShowError(c, err.Error())
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}
