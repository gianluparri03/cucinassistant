package web

import (
	"embed"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"cucinassistant/configs"
)

//go:embed all:assets
var assets embed.FS

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

// registerAssets adds all the assets to the router
func registerAssets(router *mux.Router) {
	// Registers the assets
	fs := cacheAssets(http.FileServerFS(assets))
	router.PathPrefix("/assets/").Handler(fs)

	// Registers the favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusSeeOther)
	})
}
