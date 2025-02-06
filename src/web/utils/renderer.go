package utils

import (
	"net/http"

	"cucinassistant/langs"
)

// render executes the given templates witth the given data.
// It adds the base template in automatic, looking if it's a
// normal request or a request made by htmx.
func render(c *Context, pages []string, data map[string]any) {
	// Adds the base template
	if c.h {
		pages = append([]string{"templates/base_hx"}, pages...)
	} else {
		pages = append([]string{"templates/base_full"}, pages...)
	}

	// Completes the pages names
	for i, p := range pages {
		pages[i] = "web/pages/" + p + ".html"
	}

	// Adds the isHx field to the data
	if data == nil {
		data = make(map[string]any)
	}
	data["IsHx"] = c.h

	// Executes the template
	langs.ExecuteTemplates(c.W, c.L, pages, data)
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
	tMsg := langs.Translate(c.L, msg, nil)
	render(c, []string{"templates/error"}, map[string]any{"Message": tMsg, "Path": path})
}

// Show shows a popup message to the user.
func Show(c *Context, msg string) {
	ShowAndRedirect(c, msg, "")
}
