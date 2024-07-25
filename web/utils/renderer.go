package utils

import (
	"html/template"
	"log/slog"
	"net/http"
)

// render executes the given templates witth the given data.
// It adds the base template in automatic, looking if it's a
// normal request or a request made by htmx.
func render(c Context, pages []string, data map[string]any) {
	// Adds the base template
	isHx := c.R.Header.Get("HX-Request") != ""
	if !isHx {
		pages = append([]string{"templates/base_full"}, pages...)
	} else {
		pages = append([]string{"templates/base_hx"}, pages...)
	}

	// Completes the pages names
	for i, p := range pages {
		pages[i] = "web/pages/" + p + ".html"
	}

	// Loads the templates
	tmpl, err := template.ParseFiles(pages...)
	if err != nil {
		slog.Error("while fetching page template:", "err", err, "pages", pages)
	}

	// Adds the isHx field to the data
	if data == nil {
		data = map[string]any{"IsHx": isHx}
	} else {
		data["IsHx"] = isHx
	}

	// Parses them
	if err = tmpl.Execute(c.W, data); err != nil {
		slog.Error("while executing page template:", "err", err, "pages", pages)
	}
}

// RenderPage renders a specific page, with some data.
// PageName must contain the subfolder and the basename, like
// "account/signup"
func RenderPage(c Context, pageName string, data map[string]any) {
	render(c, []string{"templates/body", pageName}, data)
}

// ShowError shows an error to the user
func ShowError(c Context, msg string) {
	c.W.WriteHeader(http.StatusBadRequest)
	render(c, []string{"templates/error"}, map[string]any{"Message": msg})
}
