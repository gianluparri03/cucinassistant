package web

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"

	"cucinassistant/configs"
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
	if err := http.ListenAndServe(":"+configs.Port, router); err != nil {
		slog.Error("while running web server:", "err", err)
	}
}

// registerAssets adds all the assets to the router
func registerAssets(router *mux.Router) {
	// Registers the assets
	fs := cacheAssets(http.FileServer(http.Dir("web/assets")))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Registers the favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/assets/logo_round.png")
	})

	// Registers the 404 handler
	router.NotFoundHandler = utils.Handler(func(c *utils.Context) error {
		utils.ShowError(c, "MSG_PAGE_NOT_FOUND", "/", http.StatusNotFound)
		return nil
	})
}

// cacheAssets adds the ETag to the resource. The ETag is set to the current
// version.
func cacheAssets(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := "\"" + strconv.Itoa(configs.VersionCode) + "\""

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", version)

		if w.Header().Get("If-None-Match") == version {
			w.WriteHeader(http.StatusNotModified)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
