package web

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log/slog"
	"net/http"
)

func renderPage(w http.ResponseWriter, r *http.Request, s *sessions.Session, page string, data map[string]any) {
	// Chooses the right base template
	var baseTemplate string
	if r.Header.Get("HX-Request") == "" {
		baseTemplate = "base"
	} else {
		baseTemplate = "base_hx"
	}

	// Loads the templates
	tmpl, err := template.ParseFiles(
		"web/pages/templates/"+baseTemplate+".html",
		"web/pages/templates/boot_script.html",
		"web/pages/"+page+".html",
	)

	if err != nil {
		slog.Error("while fetching page template:", "err", err, "template", page)
	}

	// Adds the error (if any)
	if err, found := s.Values["Error"]; found {
		if data == nil {
			data = map[string]any{"Error": err}
		} else {
			data["Error"] = err
		}

		delete(s.Values, "Error")
		saveSession(w, r, s)
	}

	// And finally parses them
	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("while executing page template:", "err", err, "template", page)
	}
}

func showError(w http.ResponseWriter, r *http.Request, s *sessions.Session, err error) {
	s.Values["Error"] = err.Error()
	saveSession(w, r, s)
	http.Redirect(w, r, r.URL.String(), 301)
}
