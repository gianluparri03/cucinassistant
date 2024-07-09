package web

import (
	"html/template"
	"log"
	"net/http"

	"cucinassistant/config"
)

func Start() {
	registerRoutes()

	// Starts the server
	if err := http.ListenAndServe(config.Runtime.ServerAddress, nil); err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func registerRoutes() {
	// Static files
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/account/accedi/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, "account/signin", map[string]any{"NavigationDisabled": true})
	})

	http.HandleFunc("/account/recupera_password/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, "account/password_recover", map[string]any{})
	})
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
