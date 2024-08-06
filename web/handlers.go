package web

import (
	"net/http"
	"strconv"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web/utils"
	"github.com/gorilla/mux"
)

// RegisterAll adds all the routes to the router
func RegisterAll(router *mux.Router) {
	// Prepares all the routes
	endpoints := []utils.Endpoint{
		{
			Path:    "/",
			GetPage: "misc/home",
			GetData: func(c utils.Context) map[string]any {
				user, _ := database.GetUser(c.UID)
				return map[string]any{"Username": user.Username}
			},
			PostDisabled: true,
		},

		{
			Path:    "/info",
			GetPage: "misc/info",
			GetData: func(c utils.Context) map[string]any {
				return map[string]any{
					"Config":      config.Runtime,
					"Version":     config.Version,
					"UsersNumber": database.GetUsersNumber(),
				}
			},
			PostDisabled: true,
		},

		{
			Path:    "/shopping_list",
			GetPage: "shopping_list/view",
			GetData: func(c utils.Context) map[string]any {
				user, _ := database.GetUser(c.UID)
				list, _ := user.GetShoppingList()
				return map[string]any{"List": list}
			},
			PostDisabled: true,
		},
		{
			Path:        "/shopping_list/append",
			GetPage:     "shopping_list/append",
			PostHandler: appendEntries,
		},
		{
			Path:        "/shopping_list/{EID}/toggle",
			GetDisabled: true,
			PostHandler: toggleEntry,
		},
		{
			Path:    "/shopping_list/{EID}/edit",
			GetPage: "shopping_list/edit",
			GetData: func(c utils.Context) map[string]any {
				EID, _ := strconv.Atoi(mux.Vars(c.R)["EID"])
				user, _ := database.GetUser(c.UID)
				entry, _ := user.GetEntry(EID)
				return map[string]any{"Name": entry.Name}
			},
			PostHandler: editEntry,
		},
		{
			Path:        "/shopping_list/clear",
			GetDisabled: true,
			PostHandler: clearEntries,
		},

		{
			Path:        "/user/signup",
			Unprotected: true,
			GetPage:     "user/signup",
			PostHandler: signUp,
		},
		{
			Path:        "/user/signin",
			Unprotected: true,
			GetPage:     "user/signin",
			PostHandler: signIn,
		},
		{
			Path:        "/user/signout",
			GetDisabled: true,
			PostHandler: signOut,
		},
		{
			Path:        "/user/forgot_password",
			GetPage:     "user/forgot_password",
			PostHandler: forgotPassword,
		},
		{
			Path:    "/user/reset_password",
			GetPage: "user/reset_password",
			GetData: func(c utils.Context) map[string]any {
				return map[string]any{"Token": c.R.URL.Query().Get("token")}
			},
			PostHandler: resetPassword,
		},
		{
			Path:         "/user/settings",
			GetPage:      "user/settings",
			PostDisabled: true,
		},
		{
			Path:    "/user/change_username",
			GetPage: "user/change_username",
			GetData: func(c utils.Context) map[string]any {
				user, _ := database.GetUser(c.UID)
				return map[string]any{"Username": user.Username}
			},
			PostHandler: changeUsername,
		},
		{
			Path:    "/user/change_email",
			GetPage: "user/change_email",
			GetData: func(c utils.Context) map[string]any {
				user, _ := database.GetUser(c.UID)
				return map[string]any{"Email": user.Email}
			},
			PostHandler: changeEmail,
		},
		{
			Path:        "/user/change_password",
			GetPage:     "user/change_password",
			PostHandler: changePassword,
		},
		{
			Path:    "/user/delete_1",
			GetPage: "user/delete",
			GetData: func(c utils.Context) map[string]any {
				return map[string]any{"Warning": true}
			},
			PostHandler: deleteUser1,
		},
		{
			Path:    "/user/delete_2",
			GetPage: "user/delete",
			GetData: func(c utils.Context) map[string]any {
				return map[string]any{"Token": c.R.URL.Query().Get("token")}
			},
			PostHandler: deleteUser2,
		},
	}

	// Registers the assets
	fs := http.FileServer(http.Dir("web/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Registers the favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusMovedPermanently)
	})

	// Registers all the endpoints
	for _, e := range endpoints {
		e.Register(router)
	}

	// Registers the 404 handler
	router.NotFoundHandler = utils.Handler(func(c utils.Context) {
		utils.ShowError(c, "Pagina non trovata", "/")
	})
}
