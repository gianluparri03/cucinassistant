package web

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"

	"cucinassistant/config"
	"cucinassistant/web/utils"
)

// Start creates and starts the web server
func Start() {
	// Creates the router
	router := mux.NewRouter()

	// Registers all the endpoints
	for _, e := range endpoints {
		e.Register(router)
	}

	// Registers all the assets
	RegisterAssets(router)

	// Prepares the session storage
	utils.InitSessionStore()

	// Starts the server
	if err := http.ListenAndServe(config.Runtime.ServerAddress, router); err != nil {
		slog.Error("while running web server:", "err", err)
	}
}

// RegisterAssets adds all the assets to the router
func RegisterAssets(router *mux.Router) {
	// Registers the assets
	fs := http.FileServer(http.Dir("web/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Registers the favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusMovedPermanently)
	})

	// Registers the 404 handler
	router.NotFoundHandler = utils.Handler(func(c utils.Context) {
		utils.ShowMessage(c, "Pagina non trovata", "/")
	})
}
