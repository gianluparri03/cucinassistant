package web

import (
	"github.com/gorilla/mux"
	"github.com/srinathgs/mysqlstore"
	"html/template"
	"log"
	"net/http"

	"cucinassistant/config"
	"cucinassistant/database"
)

var store *mysqlstore.MySQLStore

func Start() {
	// Creates the router
	router := createRouter()

	// Prepares the session storage
	var err error
	store, err = mysqlstore.NewMySQLStoreFromConnection(
		database.DB,
		"sessions",
		"/",
		60*60*24*90,
	)
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}

	// Starts the server
	if err := http.ListenAndServe(config.Runtime.ServerAddress, router); err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func createRouter() (router *mux.Router) {
	router = mux.NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("web/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Favicon
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/logo_round.png", http.StatusMovedPermanently)
	})

	// TODO move away
	router.HandleFunc("/account/accedi/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, "account/signin", map[string]any{"NavigationDisabled": true})
	})

	// TODO move away
	router.HandleFunc("/account/recupera_password/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, "account/password_recover", map[string]any{})
	})

	return
}

func renderPage(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	// Chooses the right base template
	var baseTemplate string
	if r.Header.Get("HX-Request") == "" {
		baseTemplate = "base"
	} else {
		baseTemplate = "base_hx"
	}

	// Loads the templates
	tmpl := template.Must(template.ParseFiles(
		"web/pages/templates/"+baseTemplate+".html",
		"web/pages/templates/boot_script.html",
		"web/pages/"+page+".html",
	))

	// And finally parses them
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("ERR: " + err.Error())
	}
}
