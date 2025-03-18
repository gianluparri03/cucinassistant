package utils

import (
	"github.com/a-h/templ"
	"net/http"

	"cucinassistant/langs"
	"cucinassistant/web/components"
)

// render renders the given page.
// If it's an htmx request it will render only content; otherwise it will
// render the complete structure, with the body and the message.
func render(c *Context, body, message, content templ.Component) {
	if !c.h {
		content = components.TemplateBase(c.L, body, message)
	}

	content.Render(langs.Lang(c.L), c.W)
}

// RenderComponent renders a component
func RenderComponent(c *Context, page templ.Component) {
	render(c, page, components.TemplateEmpty(), page)
}

// ShowError shows a popup error to the user, returning the
// given http status code.
// If path is set, it will redirects it to the given path
func ShowError(c *Context, msg string, path string, code int) {
	c.W.Header().Set("HX-Retarget", "#message-container")
	c.W.WriteHeader(code)
	page := components.TemplateMessage(msg, path, c.h)
	render(c, components.TemplateEmpty(), page, page)
}

// ShowMessage shows a popup message to the user.
// If path is set, it will redirects it to the given path
func ShowMessage(c *Context, msg string, path string) {
	ShowError(c, msg, path, http.StatusOK)
}

// Redirect redirects to a given path
func Redirect(c *Context, path string) {
	http.Redirect(c.W, c.R, path, http.StatusSeeOther)
}
