package utils

import (
	"net/http"
)

// Redirect redirects to a given path
func Redirect(c Context, path string) {
	http.Redirect(c.W, c.R, path, http.StatusSeeOther)
}
