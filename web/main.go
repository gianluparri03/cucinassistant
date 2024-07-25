package web

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"

	"cucinassistant/config"
	"cucinassistant/web/handlers"
	"cucinassistant/web/utils"
)

// Start creates and starts the web server
func Start() {
	// Creates the router
	router := mux.NewRouter()

	// Registers all the handlers
	handlers.RegisterAssetsHandlers(router)
	handlers.RegisterAccountHandlers(router)
	handlers.RegisterMiscHandlers(router)

	// Prepares the session storage
	utils.InitSessionStore()

	// Starts the server
	if err := http.ListenAndServe(config.Runtime.ServerAddress, router); err != nil {
		slog.Error("while running web server:", "err", err)
	}
}