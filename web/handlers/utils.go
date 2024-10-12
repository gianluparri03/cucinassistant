package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/web/utils"
)

var (
	MSG_EMAIL_SENT = "Ti abbiamo inviato una mail. Controlla la casella di posta."
)

// getID is used to retrieve and ID from the URL.
// The third parameter is the error returned if
// somethign goes wrong.
func getID(c *utils.Context, name string, notFound error) (int, error) {
	ID, err := strconv.Atoi(mux.Vars(c.R)[name])
	if err != nil {
		return 0, notFound
	} else {
		return ID, nil
	}
}
