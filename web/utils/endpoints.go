package utils

import (
	"github.com/gorilla/mux"
)

// Endpoint contains all the info to register automatically
// the GET and POST route for an endpoint
type Endpoint struct {
	// Path is the endpoint's path
	Path string

	// Unprotected indicates whether the user has to be logged
	// in to use this endpoint (bot for GET and POST requests)
	Unprotected bool

	// PostHandler is the function executed on GET requests.
	// If not set, the endpoint will show an error on GET requests
	GetHandler func(c Context)

	// PostHandler is the function executed on POST requests
	// If not set, the endpoint will show an error on POST requests
	PostHandler func(c Context)
}

// Register adds the endpoint to the router
func (e Endpoint) Register(router *mux.Router) {
	// Returns an error message if the method
	// is not supported
	unknownHandler := func(c Context) {
		ShowMessage(c, "Richiesta sconosciuta", "/")
	}

	// Prepares the GET handler
	var get func(c Context)
	if e.GetHandler != nil {
		get = e.GetHandler
	} else {
		get = unknownHandler
	}

	// Prepares the POST handler
	var post func(c Context)
	if e.PostHandler != nil {
		post = e.PostHandler
	} else {
		post = unknownHandler
	}

	// Registers them
	if e.Unprotected {
		router.Handle(e.Path, Handler(get)).Methods("GET")
		router.Handle(e.Path, Handler(post)).Methods("POST")
	} else {
		router.Handle(e.Path, PHandler(get)).Methods("GET")
		router.Handle(e.Path, PHandler(post)).Methods("POST")
	}
}
