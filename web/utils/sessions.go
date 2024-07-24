package utils

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"strings"

	"cucinassistant/config"
)

var store *sessions.CookieStore

// InitSessionStore initializes the cookie session store.
// It lasts 90 days.
func InitSessionStore() {
	// Initializes the session store
	store = sessions.NewCookieStore([]byte(config.Runtime.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 90,
		Secure:   strings.HasPrefix(config.Runtime.BaseURL, "https://"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// SaveSession saves the session content
func SaveSession(c Context) {
	// Saves the session
	if err := c.S.Save(c.R, c.W); err != nil {
		slog.Warn("during session saving:", "err", err)
	}
}
