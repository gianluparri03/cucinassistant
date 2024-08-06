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

	// GetDisabled indicates whether to return the page or an error
	GetDisabled bool

	// GetPage is the page to be rendered for GET requests.
	// See RenderPage for more info.
	GetPage string

	// GetData is a function that calculates all the auxiliary data
	// needed to render the page correctly.
	// See RenderPage for more info.
	GetData func(c Context) map[string]any

	// PostDisabled indicates whether to use PostHandler or
	// show an error on POST requests
	PostDisabled bool

	// PostHandler is the function executed on POST requests
	// if PostDisabled is set to true
	PostHandler func(c Context)
}

// Register adds the endpoint to the router
func (e Endpoint) Register(router *mux.Router) {
	// Prepares the GET handler
	GetHandler := func(c Context) {
		if e.GetDisabled {
			ShowError(c, "Richiesta sconosciuta", "/")
		} else {
			if e.GetData != nil {
				RenderPage(c, e.GetPage, e.GetData(c))
			} else {
				RenderPage(c, e.GetPage, nil)
			}
		}
	}

	// Prepares the POST handler
	PostHandler := func(c Context) {
		if e.PostDisabled {
			ShowError(c, "Richiesta sconosciuta", "/")
		} else {
			e.PostHandler(c)
		}
	}

	// Registers them
	if e.Unprotected {
		router.Handle(e.Path, Handler(GetHandler)).Methods("GET")
		router.Handle(e.Path, Handler(PostHandler)).Methods("POST")
	} else {
		router.Handle(e.Path, PHandler(GetHandler)).Methods("GET")
		router.Handle(e.Path, PHandler(PostHandler)).Methods("POST")
	}
}
