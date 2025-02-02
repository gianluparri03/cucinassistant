package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/langs"
	"cucinassistant/web/utils"
)

// GetLang renders /lang
// If it has the 'tag' query param, it sets the current
// language; otherwise, it shows the form
func GetLang(c *utils.Context) error {
	if tag := c.R.URL.Query().Get("tag"); tag == "" {
		utils.RenderPage(c, "misc/lang", map[string]any{
			"Langs":   langs.Available,
			"Current": c.L,
		})
	} else {
		utils.SetLang(c, tag)
	}

	return nil
}

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
