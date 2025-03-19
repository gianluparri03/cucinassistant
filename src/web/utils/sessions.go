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
	"cucinassistant/langs"
)

// store is used to store the cookies
var store *pgstore.PGStore

// sessionName is the session's name
var sessionName = "ca_session"

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
func SaveUID(c *Context, UID int, msg langs.String) {
	c.s.Values["UID"] = UID
	if err := c.s.Save(c.R, c.W); err != nil {
		slog.Error("while saving session:", "err", err)
	}

	if msg != langs.STR_NONE {
		ShowMessage(c, msg, "/")
	} else {
		Redirect(c, "/")
	}
}

// DropUID drops the UID from the session.
// It also redirects to /user/signin, with an optional message
func DropUID(c *Context, msg langs.String) {
	delete(c.s.Values, "UID")
	if err := c.s.Save(c.R, c.W); err != nil {
		slog.Error("while saving session:", "err", err)
	}

	if msg != langs.STR_NONE {
		ShowMessage(c, msg, "/user/signin")
	} else {
		Redirect(c, "/user/signin")
	}
}

// SetLang sets the session language
func SetLang(c *Context, lang string) {
	if _, found := langs.Available[lang]; !found {
		ShowError(c, langs.STR_UNKNOWN_LANG, "", http.StatusBadRequest)
		return
	}

	c.s.Values["Lang"] = lang
	c.L = lang
	c.s.Save(c.R, c.W)

	if c.h {
		ShowMessage(c, langs.STR_LANG_CHANGED, "/")
	} else {
		Redirect(c, "/")
	}
}
