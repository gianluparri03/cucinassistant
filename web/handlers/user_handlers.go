package handlers

import (
	"github.com/gorilla/mux"

	"cucinassistant/web/utils"
)

// RegisterUserHandlers registers all the handlers for the
// endpoints beginning with /user/
func RegisterUserHandlers(router *mux.Router) {
	router.Handle("/user", utils.Handler(handleGetUser)).Methods("GET")
	router.Handle("/user", utils.Handler(handlePostUser)).Methods("POST")
	router.Handle("/user/settings", utils.Handler(handleGetUserSettings)).Methods("GET")
}

// handleGetUser is called for GET requests at /user.
// It looks at the query string for "action" to decide which form to
// render. The possible values are "signup", "signin", "forgot_password"
// and "reset_password".
// In all the other cases, it renders an error.
func handleGetUser(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	switch c.R.URL.Query().Get("action") {
	case "signup":
		utils.RenderPage(c, "user/signup", nil)
	case "signin":
		utils.RenderPage(c, "user/signin", nil)
	case "forgot_password":
		utils.RenderPage(c, "user/forgot_password", nil)
	case "reset_password":
		utils.RenderPage(c, "user/reset_password", map[string]any{"Token": c.R.URL.Query().Get("token")})
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}

// handlePostUser is called for POST requests at /user
// It looks at the request value for "action" to decide what to do.
// The possible values are "signup", "signin", "signout", "forgot_password"
// and "reset_password".
// In all the other cases, it returns an error.
func handlePostUser(c utils.Context) {
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

// handleGetUserSettings is called for GET requests at /user/settings.
// If there isn't a query string for "action", it renders the settings dashboard.
// In all the other cases, it renders an error.
func handleGetUserSettings(c utils.Context) {
	// Decides which page to render based on the value
	// of the action field in the post body
	switch c.R.URL.Query().Get("action") {
	case "":
		utils.RenderPage(c, "user/settings", nil)
	default:
		utils.ShowError(c, "Richiesta sconosciuta")
	}
}
