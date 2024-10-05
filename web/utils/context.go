package utils

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"

	"cucinassistant/database"
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

	// User is the user currently logged in
	U *database.User
}

// Handler is a function that accepts a context and returns an error.
// It is used to serve http requests. If an error is returned, it will
// be showed to the user.
type Handler func(*Context) error

// ServeHTTP is used by net/http to serve an http request
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "] " + r.URL.String())

	// Executes the handler
	s, err := store.Get(r, sessionName)
	if err != nil {
		slog.Warn("while retrieving session:", "err", err)
	}

	c := &Context{W: w, R: r, S: s}
	err = h(c)

	// Shows the error (if present)
	if err != nil {
		Show(c, err.Error())
	}
}

// PHandler (ProtectedHandler) is like an Handler,
// but redirects unregistered users away.
// It also provides the user in the context.
type PHandler func(*Context) error

// ServeHTTP is used by net/http to serve an http request
func (ph PHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "*] " + r.URL.String())

	// Lets the handlers use the session
	s, err := store.Get(r, sessionName)
	if err != nil {
		slog.Warn("while retrieving session:", "err", err)
	}

	c := &Context{W: w, R: r, S: s}

	// Gets the UID from the cookies
	if rawUID, found := s.Values["UID"]; !found {
		// If there isn't an UID, redirects to the signin
		Redirect(c, "/user/signin")
	} else {
		// Fetches the user from the database
		var err error
		if c.U, err = database.GetUser("UID", rawUID.(int)); err == nil {
			// If all was okay, executes the handler
			err = ph(c)
		}

		// Shows the error (if present)
		if err != nil {
			if err == database.ERR_USER_UNKNOWN {
				DropUID(c, "")
			}

			Show(c, err.Error())
		}
	}
}
