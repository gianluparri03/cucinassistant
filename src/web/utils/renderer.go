package utils

import (
	"github.com/a-h/templ"
	"net/http"

	"cucinassistant/langs"
	"cucinassistant/web/pages"
)

// render renders the given page.
// If it's an htmx request it will render only content; otherwise it will
// render the complete structure, with the body and the message.
func render(c *Context, body, message, content templ.Component) {
	if !c.h {
		content = pages.TemplateBase(c.L, body, message)
	}

	content.Render(langs.Lang(c.L), c.W)
}

// RenderPage renders a page (a body)
func RenderPage(c *Context, page templ.Component) {
	render(c, page, pages.TemplateEmpty(), page)
}

// ShowAndRedirect shows a popup message to the user,
// then redirects it away
func ShowAndRedirect(c *Context, msg string, path string) {
	c.W.WriteHeader(http.StatusBadRequest)
	page := pages.TemplateMessage(msg, path, c.h)
	render(c, pages.TemplateEmpty(), page, page)
}

// Show shows a popup message to the user.
func Show(c *Context, msg string) {
	ShowAndRedirect(c, msg, "")
}

// Redirect redirects to a given path
func Redirect(c *Context, path string) {
	http.Redirect(c.W, c.R, path, http.StatusSeeOther)
}
