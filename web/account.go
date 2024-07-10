package web

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"

	"cucinassistant/database"
	"cucinassistant/email"
)

func registerAccountRoutes(router *mux.Router) {
	// Registers the rooutes starting with /account
	router.Handle("/account", withSession(handleGetAccount)).Methods("GET")
	router.Handle("/account", withSession(handlePostAccount)).Methods("POST")
}

func handleGetAccount(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	// Decides which page to render based on the value
	// of the action field in the query string
	if action := r.URL.Query().Get("action"); action == "signup" {
		renderPage(w, r, s, "account/signup", map[string]any{"NavigationDisabled": true})
	} else {
		http.Error(w, "Unrecognized", 500)
	}
}

func handlePostAccount(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	// Decides which page to render based on the value
	// of the action field in the request body
	if action := r.URL.Query().Get("action"); action == "signup" {
		signUpUser(w, r, s)
	} else {
		http.Error(w, "Unrecognized", 500)
	}
}

func signUpUser(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	// Fetches data
	user := &database.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Checks if all the required data is present;
	// then tries to sign it up
	var err error
	if user.Username == "" {
		err = errors.New("Nome utente mancante")
	} else if user.Email == "" {
		err = errors.New("Email mancante")
	} else if user.Password == "" {
		err = errors.New("Password mancante")
	} else if user.Password != r.FormValue("password2") {
		err = errors.New("Le due password non corrispondono")
	} else {
		err = user.SignUp()
	}

	// Ensures there's been no errors
	if err != nil {
		showError(w, r, s, err)
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	s.Values["UID"] = user.UID
	saveSession(w, r, s)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
