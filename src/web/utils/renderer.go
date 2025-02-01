package utils

import (
	"html/template"
	"log/slog"
	"net/http"
	"path"

	"cucinassistant/langs"
)

// getLang returns the desired language
func getLang(c *Context) string {
	return "it"
}

// render executes the given templates witth the given data.
// It adds the base template in automatic, looking if it's a
// normal request or a request made by htmx.
func render(c *Context, pages []string, data map[string]any) {
	// Adds the base template
	isHx := c.R.Header.Get("HX-Request") != ""
	if !isHx {
		pages = append([]string{"templates/base_full"}, pages...)
	} else {
		pages = append([]string{"templates/base_hx"}, pages...)
	}

	// Stores the template name
	name := path.Base(pages[0] + ".html")

	// Completes the pages names
	for i, p := range pages {
		pages[i] = "web/pages/" + p + ".html"
	}

	// Prepares the FuncMap
	funcs := template.FuncMap{
		"t": func(id string, data any) template.HTML {
			return template.HTML(langs.Translate(getLang(c), id, data))
		},
	}

	// Loads the templates
	tmpl, err := template.New(name).Funcs(funcs).ParseFiles(pages...)
	if err != nil {
		slog.Error("while fetching page template:", "err", err, "pages", pages)
		return
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
// "user/signup"
func RenderPage(c *Context, pageName string, data map[string]any) {
	render(c, []string{"templates/body", pageName}, data)
}

// Redirect redirects to a given path
func Redirect(c *Context, path string) {
	http.Redirect(c.W, c.R, path, http.StatusSeeOther)
}

// ShowAndRedirect shows a popup message to the user,
// then redirects him away
func ShowAndRedirect(c *Context, msg string, path string) {
	c.W.WriteHeader(http.StatusBadRequest)
	tMsg := langs.Translate(getLang(c), msg, nil)
	render(c, []string{"templates/error"}, map[string]any{"Message": tMsg, "Path": path})
}

// Show shows a popup message to the user.
func Show(c *Context, msg string) {
	ShowAndRedirect(c, msg, "")
}
