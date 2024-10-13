package handlers

import (
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// getSID returns the SID written in the url
func getSID(c *utils.Context) (SID int, err error) {
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
		utils.RenderPage(c, "storage/dashboard", map[string]any{"List": list})
	}

	return
}

// GetNewSection renders /storage/new
func GetNewSection(c *utils.Context) (err error) {
	utils.RenderPage(c, "storage/new_section", nil)
	return
}

// PostNewSection tries to create a new section
func PostNewSection(c *utils.Context) (err error) {
	var s database.Section
	if s, err = c.U.Storage().NewSection(c.R.FormValue("name")); err == nil {
		utils.ShowAndRedirect(c, "Sezione creata con successo", "/storage/"+strconv.Itoa(s.SID))
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
			utils.RenderPage(c, "storage/view_articles", map[string]any{
				"Section": section,
				"Search":  search,
			})
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
			utils.RenderPage(c, "storage/edit_section", map[string]any{"Section": section})
		}
	}

	return
}

// PostEditSection tries to change a section's name
func PostEditSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().EditSection(SID, c.R.FormValue("name")); err == nil {
			utils.ShowAndRedirect(c, "Modifiche applicate con successo", "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

// PostDeleteSection tries to delete a section
func PostDeleteSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		if err = c.U.Storage().DeleteSection(SID); err == nil {
			utils.ShowAndRedirect(c, "Sezione eliminata con successo", "/storage")
		}
	}

	return
}

// GetAddArticles renders /storage/{SID}/add
func GetAddArticles(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		var section database.Section

		if section, err = c.U.Storage().GetArticles(SID, ""); err == nil {
			utils.RenderPage(c, "storage/add_articles", map[string]any{
				"SID": section.SID,
			})
		}
	}

	return
}

// PostAddArticles tries to add articles to a section.
func PostAddArticles(c *utils.Context) (err error) {
	// Gets the destination SID
	var SID int
	if SID, err = getSID(c); err == nil {
		var articles []database.StringArticle
		c.R.ParseForm()

		// Insert all the articles that have a name
		for nameKey, nameValues := range c.R.PostForm {
			if strings.HasPrefix(nameKey, "article-") && strings.HasSuffix(nameKey, "-name") {
				if len(nameValues) > 0 && nameValues[0] != "" {
					articles = append(articles, database.StringArticle{
						Name:       nameValues[0],
						Expiration: c.R.PostFormValue(strings.Replace(nameKey, "-name", "-expiration", 1)),
						Quantity:   c.R.PostFormValue(strings.Replace(nameKey, "-name", "-quantity", 1)),
					})
				}
			}
		}

		// Tries to add them
		if err = c.U.Storage().AddArticles(SID, articles...); err == nil {
			utils.Redirect(c, "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

// GetSearchArticles renders /storage/{SID}/search
func GetSearchArticles(c *utils.Context) (err error) {
	var SID int
	if SID, err = getSID(c); err == nil {
		utils.RenderPage(c, "storage/search_articles", map[string]any{
			"SID": SID,
		})
	}

	return
}

// GetEditArticle renders /storage/{SID}/{AID}
func GetEditArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		var article database.OrderedArticle

		if article, err = c.U.Storage().GetOrderedArticle(SID, AID); err == nil {
			utils.RenderPage(c, "storage/edit_article", map[string]any{
				"SID":     SID,
				"Article": article,
			})
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

		if err = c.U.Storage().EditArticle(SID, AID, newData); err == nil {
			utils.ShowAndRedirect(c, "Modifiche salvate", "/storage/"+strconv.Itoa(SID)+"/"+strconv.Itoa(AID))
		}
	}

	return
}

// PostDeleteArticle tries to delete an article
func PostDeleteArticle(c *utils.Context) (err error) {
	var SID, AID int

	if SID, AID, err = getAID(c); err == nil {
		var next *int

		if err, next = c.U.Storage().DeleteArticle(SID, AID); err == nil {
			if next != nil {
				utils.Redirect(c, "/storage/"+strconv.Itoa(*next))
			} else {
				utils.Redirect(c, "/storage")
			}
		}
	}

	return
}
