package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"cucinassistant/config"
)

func Start() {
	// Creates the router
	router := createRouter()

	// Prepares the session storage
	initStore()

	// Starts the server
	if err := http.ListenAndServe(config.Runtime.ServerAddress, router); err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func createRouter() *mux.Router {
	router := mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("web/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusMovedPermanently)
	})

	// Registers the endpoints
	registerAccountRoutes(router)

	return router
}
