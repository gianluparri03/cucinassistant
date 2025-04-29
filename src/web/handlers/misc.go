package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

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
	return
}

func GetLang(c *utils.Context) (err error) {
	if tag := c.R.URL.Query().Get("tag"); tag == "" {
		utils.RenderComponent(c, components.Lang(langs.Available, c.L))
	} else {
		utils.SetLang(c, tag, c.R.URL.Query().Has("silent"))
	}

	return
}

func GetStats(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Stats(database.GetStats()))
	return
}
