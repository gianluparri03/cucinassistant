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
	registerAssets(router)

	// Prepares the session storage
	utils.InitSessionStore()

	// Starts the server
	if err := http.ListenAndServe(":"+config.Runtime.Port, router); err != nil {
		slog.Error("while running web server:", "err", err)
	}
}

// registerAssets adds all the assets to the router
func registerAssets(router *mux.Router) {
	// Registers the assets
	fs := http.FileServer(http.Dir("web/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Registers the service worker
	router.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/assets/service_worker.js")
	})

	// Registers the favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/assets/logo_round.png")
	})

	// Registers the 404 handler
	router.NotFoundHandler = utils.Handler(func(c *utils.Context) error {
		utils.ShowAndRedirect(c, "Pagina non trovata", "/")
		return nil
	})
}
