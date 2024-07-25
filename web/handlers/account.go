package handlers

import (
	"github.com/gorilla/mux"
	"net/url"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/web/utils"
)

// RegisterAccountHandlers registers all the handlers for the
// endpoints beginning with /account/
func RegisterAccountHandlers(router *mux.Router) {
	router.Handle("/account", utils.Handler(handleGetAccount)).Methods("GET")
	router.Handle("/account", utils.Handler(handlePostAccount)).Methods("POST")
	router.Handle("/account/settings", utils.Handler(handleGetAccountSettings)).Methods("GET")
}

// handleGetAccount is called for GET requests at /account.
// It looks at the query string for "action" to decide which form to
// render. The possible values are "signup", "signin", "forgot_password"
// and "reset_password".
// In all the other cases, it renders an error.
func handleGetAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	switch c.R.URL.Query().Get("action") {
	case "signup":
		utils.RenderPage(c, "account/signup", nil)
	case "signin":
		utils.RenderPage(c, "account/signin", nil)
	case "forgot_password":
		utils.RenderPage(c, "account/forgot_password", nil)
	case "reset_password":
		utils.RenderPage(c, "account/reset_password", map[string]any{"Token": c.R.URL.Query().Get("token")})
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}

// handlePostAccount is called for POST requests at /account
// It looks at the request value for "action" to decide what to do.
// The possible values are "signup", "signin", "signout", "forgot_password"
// and "reset_password".
// In all the other cases, it returns an error.
func handlePostAccount(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the request body
	switch c.R.FormValue("action") {
	case "signup":
		signUpUser(c)
	case "signin":
		signInUser(c)
	case "signout":
		signOutUser(c)
	case "forgot_password":
		forgotPassword(c)
	case "reset_password":
		resetPassword(c)
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

// signOutUser drops an user's session
func signOutUser(c utils.Context) {
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.Redirect(c, "/account?action=signin")
}

// forgotPassword sends an email to the user to recover its password
func forgotPassword(c utils.Context) {
	// Tries to get the user's email
	if user, err := database.GetUserFromEmail(c.R.FormValue("email")); err == nil {
		// Tries to generate its token
		if token, err := user.GenerateToken(); err == nil {
			data := map[string]any{
				"Username":  user.Username,
				"ResetLink": config.Runtime.BaseURL + "/account?action=reset_password&token=" + url.QueryEscape(token),
			}

			go email.SendMail(user.Email, "Recupero password", "reset_password", data)
		}
	}

	utils.ShowError(c, "Ti abbiamo inviato un email. Controlla la casella di posta")
}

// resetPassword makes an user resets his password
func resetPassword(c utils.Context) {
	// Fetches data
	password := c.R.FormValue("password")
	user := &database.User{
		Email: c.R.FormValue("email"),
		Token: c.R.FormValue("token"),
	}

	// Ensures the two passwords are equal
	if password != c.R.FormValue("password2") {
		utils.ShowError(c, "Le due password non corrispondono")
		return
	}

	// Tries to reset the user's password
	if err := user.ResetPassword(password); err != nil {
		utils.ShowError(c, err.Error())
	} else {
		// Saves the session and then redirects it to the homepage
		c.S.Values["UID"] = user.UID
		utils.SaveSession(c)
		utils.Redirect(c, "/")
	}
}

// handleGetAccountSettings is called for GET requests at /account/settings.
// If there isn't a query string for "action", it renders the settings dashboard.
// In all the other cases, it renders an error.
func handleGetAccountSettings(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	switch c.R.URL.Query().Get("action") {
	case "":
		utils.RenderPage(c, "account/settings", nil)
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}
