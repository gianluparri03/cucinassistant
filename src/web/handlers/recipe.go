package handlers

import (
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// getRID returns the RID written in the url
func getRID(c *utils.Context) (int, error) {
	return getID(c, "RID", database.ERR_RECIPE_NOT_FOUND)
}

// GetRecipes renders /recipes
func GetRecipes(c *utils.Context) (err error) {
	// TODO
	utils.RenderPage(c, "recipe/list", map[string]any{
		"Recipes": []database.Recipe{
			{
				RID:  1,
				Name: "Focaccia",
			},
			{
				RID:  2,
				Name: "Pasta al forno",
			},
		},
	})

	return
}

// GetNewRecipe renders /recipes/new
func GetNewRecipe(c *utils.Context) (err error) {
	utils.RenderPage(c, "recipe/new", nil)
	return
}

// GetRecipe renders /recipes/{RID}
func GetRecipe(c *utils.Context) (err error) {
	// TODO
	r := database.Recipe{
		RID:         1,
		Name:        "Focaccia",
		Ingredients: "Acqua qb\nSale\nFarina 500ml",
		Notes:       "Servire calda",
		Stars:       4,
		Directions:  "Mischiare tutto",
	}

	utils.RenderPage(c, "recipe/view", map[string]any{
		"Recipe":      r,
		"FullStars":   make([]struct{}, r.Stars),
		"EmptyStars":  make([]struct{}, 5-r.Stars),
		"Ingredients": strings.Split(r.Ingredients, "\n"),
		"Directions":  strings.Split(r.Directions, "\n"),
	})

	return
}

// GetEditRecipe renders /recipes/{RID}/edit
func GetEditRecipe(c *utils.Context) (err error) {
	// TODO
	r := database.Recipe{
		RID:         1,
		Name:        "Focaccia",
		Ingredients: "Acqua qb\nSale\nFarina 500ml",
		Notes:       "Servire calda",
		Stars:       4,
		Directions:  "Mischiare tutto",
	}

	utils.RenderPage(c, "recipe/edit", map[string]any{"Recipe": r})
	return
}
