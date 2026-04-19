package utils

import (
	"fmt"
	"github.com/a-h/templ"
	"net/http"

	"cucinassistant/configs"
	"cucinassistant/langs"
	"cucinassistant/web/components"
)

// render renders the given page.
// If it's an htmx request it will render only content; otherwise it will
// render the complete structure, with the body and the message.
func render(c *Context, body, message, content templ.Component) {
	if !c.h {
		tutorial := fmt.Sprintf("%s/%d_%s.pdf", configs.TutorialsURL, configs.VersionCode, c.L)
		content = components.TemplateBase(c.U != nil, c.L, body, message, tutorial)
	}

	content.Render(langs.Get(&c.L).Ctx(), c.W)
}

// RenderComponent renders a component
func RenderComponent(c *Context, page templ.Component) {
	render(c, page, components.TemplateEmpty(), page)
}

// RenderSide shows the sidebar. Can be accessed only via htmx.
func RenderSide(c *Context, side templ.Component) {
	if !c.h {
		Redirect(c, "/")
		return
	}

	c.W.Header().Add("HX-Retarget", "#side-container")
	render(c, side, components.TemplateEmpty(), side)
}

// ShowMessage shows a popup message to the user.
// If path is set, it will redirects it to the given path
func ShowMessage(c *Context, msg langs.String, path string) {
	ShowError(c, msg, path, http.StatusCreated)
}

// ShowError is like ShowMessage, but it also sets a status code
func ShowError(c *Context, msg langs.String, path string, status int) {
	c.W.Header().Add("HX-Retarget", "#message-container")
	c.W.WriteHeader(status)
	page := components.TemplateMessage(msg, path, c.h)
	render(c, components.TemplateEmpty(), page, page)
}

// Redirect redirects to a given path
func Redirect(c *Context, path string) {
	http.Redirect(c.W, c.R, path, http.StatusSeeOther)
}
