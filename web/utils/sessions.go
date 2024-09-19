package utils

import (
	"github.com/gorilla/sessions"
	"net/http"
	"strings"

	"cucinassistant/config"
)

var store *sessions.CookieStore

// InitSessionStore initializes the cookie session store.
// It lasts 90 days.
func InitSessionStore() {
	// Initializes the session store
	store = sessions.NewCookieStore([]byte(config.Runtime.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 90,
		Secure:   strings.HasPrefix(config.Runtime.BaseURL, "https://"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// SaveUID adds the UID to the session.
// It also redirects to /, with an optional message
func SaveUID(c *Context, UID int, msg string) {
	c.S.Values["UID"] = UID
	c.S.Save(c.R, c.W)
	ShowAndRedirect(c, msg, "/")
}

// DropUID drops the UID from the session.
// It also redirects to /user/signin, with an optional message
func DropUID(c *Context, msg string) {
	delete(c.S.Values, "UID")
	c.S.Save(c.R, c.W)
	ShowAndRedirect(c, msg, "/user/signin")
}
