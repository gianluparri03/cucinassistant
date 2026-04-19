package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	if lang := mux.Vars(c.R)["lang"]; lang != "" {
		if _, found := langs.Available[lang]; found {
			c.L = lang
		} else {
			utils.ShowError(c, langs.STR_UNKNOWN_LANG, "", http.StatusNotFound)
			return
		}
	}

	utils.RenderComponent(c, components.Info(map[string]string{
		"code":     configs.SourceURL,
		"version":  configs.Version,
		"tutorial": fmt.Sprintf("%s/%d_%s.pdf", configs.TutorialsURL, configs.VersionCode, c.L),
		"support":  configs.SupportEmail,
	}))
	return
}

func GetLang(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Lang(langs.Available, c.L))
	return
}

func PostLang(c *utils.Context) (err error) {
	utils.SetLang(c, c.R.FormValue("tag"))
	return
}

func GetSide(c *utils.Context) (err error) {
	var menus []database.Menu
	var sections []database.Section
	var recipes []database.Recipe

	switch c.R.URL.Query().Get("open") {
	case "menus":
		menus, _ = c.U.Menus().GetAll()
	case "sections":
		sections, _ = c.U.Storage().GetSections()
	case "recipes":
		recipes, _ = c.U.Recipes().GetAll()
	}

	utils.RenderSide(c, components.Side(menus, sections, recipes))
	return
}

func GetStats(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.Stats(database.GetStats()))
	return
}
