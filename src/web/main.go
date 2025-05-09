package web

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"

	"cucinassistant/configs"
	"cucinassistant/langs"
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

	// Registers the 404 handler
	router.NotFoundHandler = utils.Handler(func(c *utils.Context) error {
		utils.ShowError(c, langs.STR_PAGE_NOT_FOUND, "/", http.StatusNotFound)
		return nil
	})

	// Prepares the session storage
	utils.InitSessionStore()

	// Starts the server
	if err := http.ListenAndServe(":"+configs.Port, router); err != nil {
		slog.Error("while running web server:", "err", err)
	}
}
