package utils

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"

	"cucinassistant/database"
	"cucinassistant/langs"
)

// Context is a container for all the things needed
// to make an handler work.
type Context struct {
	// W is an http.ResponseWriter
	W http.ResponseWriter

	// R is a http.Response
	R *http.Request

	// U is the user currently logged in
	U *database.User

	// L is the language
	L string

	// s is a sessions.Session
	s *sessions.Session

	// h is true if the request is from htmx
	h bool
}

// buildContext creates a Context instance
func buildContext(w http.ResponseWriter, r *http.Request, protected bool) *Context {
	// Initializes the context with W, R and isHx
	c := Context{
		W: w,
		R: r,
		h: r.Header.Get("HX-Request") != "",
	}

	// Loads the session
	var err error
	c.s, err = store.Get(r, sessionName)
	if err != nil {
		slog.Warn("while retrieving session:", "err", err)
	}

	// Loads the language
	if lang, found := c.s.Values["Lang"]; found {
		c.L = lang.(string)
	} else {
		c.L = langs.Default.Tag
	}

	// Gets the user, if it's protected
	if protected {
		if rawUID, found := c.s.Values["UID"]; found {
			u, _ := database.GetUser("UID", rawUID.(int))
			c.U = &u
		}
	}

	return &c
}

// logRoute logs every route visited
func logRoute(c *Context, attr string) {
	slog.Debug("[" + c.R.Method + attr + "] " + c.R.URL.String())
}

// Handler is a function that accepts a context and returns an error.
// It is used to serve http requests. If an error is returned, it will
// be showed to the user.
type Handler func(*Context) error

// ServeHTTP is used by net/http to serve an http request
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := buildContext(w, r, false)
	logRoute(c, "")

	if err := h(c); err != nil {
		ShowError(c, langs.ParseError(err), "", http.StatusBadRequest)
	}
}

// PHandler (ProtectedHandler) is like an Handler,
// but redirects unregistered users away.
// It also provides the user in the context.
type PHandler func(*Context) error

// ServeHTTP is used by net/http to serve an http request
func (ph PHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := buildContext(w, r, true)
	logRoute(c, "*")

	if c.U != nil {
		if err := ph(c); err != nil {
			if err == database.ERR_USER_UNKNOWN {
				DropUID(c, langs.STR_NONE)
			}

			ShowError(c, langs.ParseError(err), "", http.StatusBadRequest)
		}
	} else {
		// If there isn't an UID, redirects to the signin
		Redirect(c, "/user/signin")
	}
}
