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
}

// Handler is a function that accepts a context
type Handler func(c Context)

// ServeHTTP is used by net/http to serve an http request.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logs each visit
	slog.Debug("[" + r.Method + "] " + r.URL.String())

	// Lets the handlers use the session
	s, _ := store.Get(r, "session")
	h(Context{w, r, s})
}
