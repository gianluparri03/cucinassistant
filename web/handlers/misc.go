package handlers

import (
	"github.com/gorilla/mux"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// RegisterMiscHandlers registers all the remaining handlers
func RegisterMiscHandlers(router *mux.Router) {
	router.Handle("/", utils.PHandler(handleGetIndex)).Methods("GET")
	router.Handle("/info", utils.Handler(handleGetInfo)).Methods("GET")
	router.NotFoundHandler = utils.Handler(handleNotFound)
}

// handleGetIndex is called for GET* requests at /
func handleGetIndex(c utils.Context) {
	data := map[string]any{"Username": database.GetUserData(c.UID).Username}
	utils.RenderPage(c, "misc/home", data)
}

// handleGetInfo is called for GET requests at /info
func handleGetInfo(c utils.Context) {
	data := map[string]any{
		"Config":      config.Runtime,
		"Version":     config.Version,
		"UsersNumber": database.GetUsersNumber(),
	}

	utils.RenderPage(c, "misc/info", data)
}

// handleNotFound is used when the page has not been found
func handleNotFound(c utils.Context) {
	utils.ShowError(c, "Pagina non trovata")
}
