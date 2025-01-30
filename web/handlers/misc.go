package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetIndex renders /
func GetIndex(c *utils.Context) error {
	utils.RenderPage(c, "misc/home", map[string]any{"Username": c.U.Username})
	return nil
}

// GetInfo renders /info
func GetInfo(c *utils.Context) error {
	utils.RenderPage(c, "misc/info", map[string]any{
		"BaseURL":     configs.BaseURL,
		"VersionCode": configs.VersionCode,
		"VersionName": configs.VersionName,
	})
	return nil
}

// GetStats renders /stats
func GetStats(c *utils.Context) error {
	utils.RenderPage(c, "misc/stats", map[string]any{
		"Stats": database.GetStats(),
	})
	return nil
}

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
