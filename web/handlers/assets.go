package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterAssetsHandlers registers the handlers for all
// the static files
func RegisterAssetsHandlers(r *mux.Router) {
	// Static files
	fs := http.FileServer(http.Dir("web/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Favicon
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusMovedPermanently)
	})
}
