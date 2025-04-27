package handlers

import (
	"github.com/gorilla/mux"
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/langs"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func getAID(c *utils.Context) (int, int, error) {
	SID, errS := getSID(c)
	AID, errA := getID(c, "AID", database.ERR_ARTICLE_NOT_FOUND)

	if errS != nil {
		return SID, AID, errS
	} else if errA != nil {
		return SID, AID, errA
	} else {
		return SID, AID, nil
	}
}

func getSID(c *utils.Context) (int, error) {
	return getID(c, "SID", database.ERR_SECTION_NOT_FOUND)
}

func GetStorage(c *utils.Context) (err error) {
	var list []database.Section

	if list, err = c.U.Storage().GetSections(); err == nil {
		utils.RenderComponent(c, components.Storage(list))
	}

	return
}

func GetStorageNew(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.StorageNew())
	return
}

func PostStorageNew(c *utils.Context) (err error) {
	var s database.Section

	if s, err = c.U.Storage().NewSection(c.R.FormValue("name")); err == nil {
		utils.Redirect(c, "/storage/"+strconv.Itoa(s.SID))
	}

	return
}

func GetStorageSection(c *utils.Context) (err error) {
	var SID int
	var section database.Section
	sn := make(map[int]string)

	if SID, err = getSID(c); err == nil {
		search := c.R.URL.Query().Get("search")
		if section, err = c.U.Storage().GetArticles(SID, search); err == nil {
			if SID == 0 {
				sections, _ := c.U.Storage().GetSections()
				for _, s := range sections {
					sn[s.SID] = s.Name
				}
			}

			utils.RenderComponent(c, components.StorageSection(section, search, sn))
		}
	}

	return
}

func GetStorageSectionAdd(c *utils.Context) (err error) {
	var SID int
	var sections []database.Section

	if SID, err = getSID(c); err == nil {
		if SID == 0 {
			sections, _ = c.U.Storage().GetSections()
		}

		utils.RenderComponent(c, components.StorageSectionAdd(SID, sections))
	}

	return
}

func PostStorageSectionAdd(c *utils.Context) (err error) {
	var articles []database.StringArticle
	c.R.ParseForm()
	prefix := "article-"

	for key, values := range c.R.PostForm {
		if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, "-name") {
			if len(values) > 0 && values[0] != "" {
				id := key[len(prefix) : len(key)-len("-name")]

				name := values[0]
				exp := c.R.PostFormValue(prefix + id + "-expiration")
				qty := c.R.PostFormValue(prefix + id + "-quantity")
				sid, _ := strconv.Atoi(c.R.PostFormValue(prefix + id + "-section"))

				articles = append(articles, database.StringArticle{
					Name: name, Expiration: exp, Quantity: qty, Section: sid,
				})
			}
		}
	}

	if err = c.U.Storage().AddArticles(articles...); err == nil {
		sid, _ := mux.Vars(c.R)["SID"]
		utils.ShowMessage(c, langs.STR_ARTICLES_ADDED, "/storage/"+sid)
	}

	return
}

func PostStorageSectionDelete(c *utils.Context) (err error) {
	var SID int

	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().DeleteSection(SID); err == nil {
			utils.ShowMessage(c, langs.STR_SECTION_DELETED, "/storage")
		}
	}

	return
}

func GetStorageSectionEdit(c *utils.Context) (err error) {
	var SID int
	var section database.Section

	if SID, err = getSID(c); err == nil {
		if section, err = c.U.Storage().GetSection(SID); err == nil {
			utils.RenderComponent(c, components.StorageSectionEdit(section))
		}
	}

	return
}

func PostStorageSectionEdit(c *utils.Context) (err error) {
	var SID int

	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().EditSection(SID, c.R.FormValue("name")); err == nil {
			utils.Redirect(c, "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

func GetStorageSectionSearch(c *utils.Context) (err error) {
	var SID int

	if SID, err = getSID(c); err == nil {
		target := "/storage/" + strconv.Itoa(SID)
		utils.RenderComponent(c, components.StorageSectionSearch(target))
	}

	return
}

func GetStorageArticle(c *utils.Context) (err error) {
	var SID, AID int
	var article database.Article

	if SID, AID, err = getAID(c); err == nil {
		if article, err = c.U.Storage().GetArticle(AID); err == nil {
			prev, next := c.U.Storage().GetNeighbours(SID, AID)
			sections, _ := c.U.Storage().GetSections()
			utils.RenderComponent(c, components.StorageArticle(SID, article, prev, next, sections))
		}
	}

	return
}

func PostStorageArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		c.R.ParseForm()
		newData := database.StringArticle{
			Name:       c.R.PostFormValue("name"),
			Expiration: c.R.PostFormValue("expiration"),
			Quantity:   c.R.PostFormValue("quantity"),
		}

		prev1, next1 := c.U.Storage().GetNeighbours(SID, AID)

		if err = c.U.Storage().EditArticle(AID, newData); err == nil {
			prev2, next2 := c.U.Storage().GetNeighbours(SID, AID)

			if prev1 == prev2 && next1 == next2 {
				utils.Redirect(c, "/storage/"+strconv.Itoa(SID)+"/"+strconv.Itoa(AID))
			} else {
				utils.ShowMessage(c, langs.STR_ORDER_CHANGED, "/storage/"+strconv.Itoa(SID))
			}
		}
	}

	return
}

func PostStorageArticleDelete(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		_, next := c.U.Storage().GetNeighbours(SID, AID)

		if err = c.U.Storage().DeleteArticle(AID); err == nil {
			path := "/storage/" + strconv.Itoa(SID)
			if next != 0 {
				path += "/" + strconv.Itoa(next)
			}

			utils.Redirect(c, path)
		}
	}

	return
}
