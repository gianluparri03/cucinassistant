package utils

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
)

// Context is a container for all the things needed
// to make an handler work.
type Context struct {
	// W is an http.ResponseWriter
	W http.ResponseWriter

	// R is a http.Response
	R *http.Request

	// S is a sessions.Session
	S *sessions.Session

	// UID is the user's uid
	UID int
}

// Handler is a function that accepts a context and can serve
// an http request
type Handler func(c Context)

// ServeHTTP is used by net/http to serve an http request
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "] " + r.URL.String())

	// Executes the handler
	s, _ := store.Get(r, "session")
	h(Context{w, r, s, 0})
}

// PHandler is a function that accepts a context and can serve
// an http request, redirecting unregistered users away
type PHandler func(c Context)

// ServeHTTP is used by net/http to serve an http request
func (ph PHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "*] " + r.URL.String())

	// Lets the handlers use the session
	s, _ := store.Get(r, "session")
	c := Context{w, r, s, 0}

	// Redirects away unregistered users
	if data, found := s.Values["UID"]; found {
		var okay bool
		if c.UID, okay = data.(int); !okay {
			// If the saved is not an int, it drops it
			// and redirects the user away
			delete(s.Values, "UID")
			SaveSession(c)
			Redirect(c, "/user/signin")
			return
		}
	} else {
		// If there isn't an UID at all, it
		// redirects the user away
		Redirect(c, "/user/signin")
		return
	}

	// Executes the handler
	ph(c)
}
