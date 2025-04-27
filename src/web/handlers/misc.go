package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/langs"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func getID(c *utils.Context, name string, notFound error) (int, error) {
	ID, err := strconv.Atoi(mux.Vars(c.R)[name])
	if err != nil {
		return 0, notFound
	} else {
		return ID, nil
	}
}

func GetIndex(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Index(c.U.Username))
	return
}

func GetInfo(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Info(map[string]any{
		"BaseURL":     configs.BaseURL,
		"VersionCode": configs.VersionCode,
		"VersionName": configs.VersionName,
	}))
	return
}

func GetLang(c *utils.Context) (err error) {
	// If it has the 'tag' query param, it sets the current
	// language; otherwise, it shows the form
	if tag := c.R.URL.Query().Get("tag"); tag == "" {
		utils.RenderComponent(c, components.Lang(langs.Available, c.L))
	} else {
		utils.SetLang(c, tag)
	}

	return
}

func GetStats(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Stats(database.GetStats()))
	return
}
