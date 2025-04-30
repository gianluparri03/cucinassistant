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
