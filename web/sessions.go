package web

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"strings"

	"cucinassistant/config"
)

type withSession func(w http.ResponseWriter, r *http.Request, s *sessions.Session)

func (h withSession) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "] " + r.URL.String())

	// Lets the handlers use the session
	s, _ := store.Get(r, "session")
	h(w, r, s)
}

var store *sessions.CookieStore

func initStore() {
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

func saveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	// Saves the session
	if err := s.Save(r, w); err != nil {
		slog.Warn("during session saving:", "err", err)
	}
}
