package utils

import (
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"cucinassistant/configs"
)

// store is used to store the cookies
var store *pgstore.PGStore

// sessionName is the session's name
var sessionName = "ca_session"

// sessionField is the field in which the UID is saved
var sessionField = "UID"

// InitSessionStore initializes the cookie session store.
// It lasts 90 days.
func InitSessionStore() {
	// Initializes the session store
	var err error
	store, err = pgstore.NewPGStore(configs.Database, []byte(configs.SessionSecret))
	if err != nil {
		slog.Error("while initializing session store:", "err", err)
		os.Exit(1)
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 90,
		Secure:   strings.HasPrefix(configs.BaseURL, "https://"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	// Cleanup the store every 24 hours
	defer store.StopCleanup(store.Cleanup(time.Hour * 24))
}

// SaveUID adds the UID to the session.
// It also redirects to /, with an optional message
func SaveUID(c *Context, UID int, msg string) {
	c.S.Values[sessionField] = UID
	if err := c.S.Save(c.R, c.W); err != nil {
		slog.Error("while saving session:", "err", err)
	}

	ShowAndRedirect(c, msg, "/")
}

// DropUID drops the UID from the session.
// It also redirects to /user/signin, with an optional message
func DropUID(c *Context, msg string) {
	delete(c.S.Values, sessionField)
	if err := c.S.Save(c.R, c.W); err != nil {
		slog.Error("while saving session:", "err", err)
	}

	ShowAndRedirect(c, msg, "/user/signin")
}
