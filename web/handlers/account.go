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
	// Registers the rooutes starting with /account
	router.Handle("/account", utils.Handler(handleGetAccount)).Methods("GET")
	router.Handle("/account", utils.Handler(handlePostAccount)).Methods("POST")
}

// handleGetAccount is called for GET requests at /account.
// If the query string for "action" is "signup", it renders the signup form.
// In all the other cases, it renders an error.
func handleGetAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	if action := c.R.URL.Query().Get("action"); action == "signup" {
		utils.RenderPage(c, "account/signup", map[string]any{"NavigationDisabled": true})
	} else {
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}

// handlePostAccount is called for POST requests at /account
// If the request value for "action" is "signup", it tries to sign up the user.
// In all the other cases, it returns an error.
func handlePostAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the request body
	if action := c.R.FormValue("action"); action == "signup" {
		signUpUser(c)
	} else {
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

	// Checks if all the required data is present;
	// then tries to sign it up
	var err string
	if user.Username == "" {
		err = "Nome utente mancante"
	} else if user.Email == "" {
		err = "Email mancante"
	} else if user.Password == "" {
		err = "Password mancante"
	} else if user.Password != c.R.FormValue("password2") {
		err = "Le due password non corrispondono"
    } else if err_ := user.SignUp(); err_ != nil {
		err = err_.Error()
	}

	// Ensures there's been no errors
	if err != "" {
		utils.ShowError(c, err)
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}
