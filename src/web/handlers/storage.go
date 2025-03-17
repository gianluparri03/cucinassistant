package handlers

import (
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/pages"
	"cucinassistant/web/utils"
)

// getSID returns the SID written in the url
func getSID(c *utils.Context) (int, error) {
	return getID(c, "SID", database.ERR_SECTION_NOT_FOUND)
}

// getAID returns the SID and the AID written in the url
func getAID(c *utils.Context) (SID int, AID int, err error) {
	SID, err = getSID(c)
	if err == nil {
		AID, err = getID(c, "AID", database.ERR_ARTICLE_NOT_FOUND)
	}

	return
}

// GetSections renders /storage
func GetSections(c *utils.Context) (err error) {
	var list []database.Section
	if list, err = c.U.Storage().GetSections(); err == nil {
		utils.RenderPage(c, pages.StorageDashboard(list))
	}

	return
}

// GetNewSection renders /storage/new
func GetNewSection(c *utils.Context) (err error) {
	utils.RenderPage(c, pages.StorageNewSection())
	return
}

// PostNewSection tries to create a new section
func PostNewSection(c *utils.Context) (err error) {
	var s database.Section
	if s, err = c.U.Storage().NewSection(c.R.FormValue("name")); err == nil {
		utils.Redirect(c, "/storage/"+strconv.Itoa(s.SID))
	}

	return
}

// GetArticles renders /storage/{SID}
func GetArticles(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		search := c.R.URL.Query().Get("search")

		var section database.Section
		if section, err = c.U.Storage().GetArticles(SID, search); err == nil {
			utils.RenderPage(c, pages.StorageViewArticles(section, search))
		}
	}

	return
}

// GetEditSection renders /storage/{SID}/edit
func GetEditSection(c *utils.Context) (err error) {
	// Retrieves the SID
	var SID int
	if SID, err = getSID(c); err == nil {
		var section database.Section
		if section, err = c.U.Storage().GetSection(SID); err == nil {
			utils.RenderPage(c, pages.StorageEditSection(section))
		}
	}

	return
}

// PostEditSection tries to change a section's name
func PostEditSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().EditSection(SID, c.R.FormValue("name")); err == nil {
			utils.Redirect(c, "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

// PostDeleteSection tries to delete a section
func PostDeleteSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().DeleteSection(SID); err == nil {
			utils.ShowAndRedirect(c, "MSG_SECTION_DELETED", "/storage")
		}
	}

	return
}

// GetAddArticlesSection renders /storage/{SID}/add
func GetAddArticlesSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		var section database.Section
		if section, err = c.U.Storage().GetArticles(SID, ""); err == nil {
			utils.RenderPage(c, pages.StorageAddArticles(section.SID, []database.Section{}))
		}
	}

	return
}

// PostAddArticlesSection tries to add articles to a section.
func PostAddArticlesSection(c *utils.Context) (err error) {
	// Gets the destination SID
	var SID int
	if SID, err = getSID(c); err == nil {
		var articles []database.StringArticle
		c.R.ParseForm()
		prefix := "article-"

		// Parse all the articles that have a name
		for key, values := range c.R.PostForm {
			if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, "-name") {
				if len(values) > 0 && values[0] != "" {
					id := key[len(prefix) : len(key)-len("-name")]

					name := values[0]
					exp := c.R.PostFormValue(prefix + id + "-expiration")
					qty := c.R.PostFormValue(prefix + id + "-quantity")

					articles = append(articles, database.StringArticle{
						Name: name, Expiration: exp, Quantity: qty, Section: SID,
					})
				}
			}
		}

		// Tries to add them
		if err = c.U.Storage().AddArticles(articles...); err == nil {
			utils.Redirect(c, "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

// GetAddArticlesCommon renders /storage/add
func GetAddArticlesCommon(c *utils.Context) (err error) {
	var sections []database.Section

	if sections, err = c.U.Storage().GetSections(); err == nil {
		utils.RenderPage(c, pages.StorageAddArticles(0, sections))
	}

	return
}

// PostAddArticlesCommon tries to add articles to multiple sections
func PostAddArticlesCommon(c *utils.Context) (err error) {
	var articles []database.StringArticle
	c.R.ParseForm()
	prefix := "article-"

	// Parse all the articles that have a name
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

	// Tries to add them
	if err = c.U.Storage().AddArticles(articles...); err == nil {
		utils.ShowAndRedirect(c, "MSG_ARTICLES_ADDED", "/storage")
	}

	return
}

// GetSearchArticles renders /storage/{SID}/search
func GetSearchArticles(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		utils.RenderPage(c, pages.StorageSearchArticles(SID))
	}

	return
}

// GetEditArticle renders /storage/{SID}/{AID}
func GetEditArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		var article database.Article
		if article, err = c.U.Storage().GetArticle(AID); err == nil {
			utils.RenderPage(c, pages.StorageEditArticle(SID, article))
		}
	}

	return
}

// PostEditArticle tries to edit an article
func PostEditArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		c.R.ParseForm()
		newData := database.StringArticle{
			Name:       c.R.PostFormValue("name"),
			Expiration: c.R.PostFormValue("expiration"),
			Quantity:   c.R.PostFormValue("quantity"),
		}

		var changed bool
		if err, changed = c.U.Storage().EditArticle(AID, newData); err == nil {
			if !changed {
				utils.Redirect(c, "/storage/"+strconv.Itoa(SID)+"/"+strconv.Itoa(AID))
			} else {
				utils.ShowAndRedirect(c, "MSG_ORDER_CHANGED", "/storage/"+strconv.Itoa(SID))
			}
		}
	}

	return
}

// PostDeleteArticle tries to delete an article
func PostDeleteArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		var next *int

		if err, next = c.U.Storage().DeleteArticle(AID); err == nil {
			path := "/storage/" + strconv.Itoa(SID)
			if next != nil {
				path += "/" + strconv.Itoa(*next)
			}

			utils.Redirect(c, path)
		}
	}

	return
}
