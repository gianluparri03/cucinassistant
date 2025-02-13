package handlers

import (
	"strconv"
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
	var recipes []database.Recipe

	if recipes, err = c.U.Recipes().GetAll(); err == nil {
		utils.RenderPage(c, "recipe/list", map[string]any{
			"Recipes": recipes,
		})
	}

	return
}

// GetNewRecipe renders /recipes/new
func GetNewRecipe(c *utils.Context) (err error) {
	utils.RenderPage(c, "recipe/new", nil)
	return
}

// PostNewRecipe tries to create a new recipe
func PostNewRecipe(c *utils.Context) (err error) {
	var recipe database.Recipe
	if recipe, err = c.U.Recipes().New(c.R.FormValue("name")); err == nil {
		utils.Redirect(c, "/recipes/"+strconv.Itoa(recipe.RID)+"/edit")
	}

	return
}

// GetRecipe renders /recipes/{RID}
func GetRecipe(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		var recipe database.Recipe

		if recipe, err = c.U.Recipes().GetOne(RID); err == nil {
			utils.RenderPage(c, "recipe/view", map[string]any{
				"Recipe":      recipe,
				"FullStars":   make([]struct{}, recipe.Stars/2),
				"HalfStars":   make([]struct{}, recipe.Stars%2),
				"EmptyStars":  make([]struct{}, (10-recipe.Stars)/2),
				"Ingredients": strings.Split(recipe.Ingredients, "\n"),
				"Directions":  strings.Split(recipe.Directions, "\n"),
			})
		}
	}

	return
}

// GetEditRecipe renders /recipes/{RID}/edit
func GetEditRecipe(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		var recipe database.Recipe
		if recipe, err = c.U.Recipes().GetOne(RID); err == nil {
			utils.RenderPage(c, "recipe/edit", map[string]any{
				"Recipe": recipe,
				"Stars":  float32(recipe.Stars) / 2,
			})
		}
	}

	return
}

// PostEditRecipe tries to edit a recipe
func PostEditRecipe(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		// Fetches data
		stars, _ := strconv.ParseFloat(c.R.FormValue("stars"), 32)
		newData := database.Recipe{
			Name:        c.R.FormValue("name"),
			Stars:       int(stars * 2),
			Ingredients: c.R.FormValue("ingredients"),
			Directions:  c.R.FormValue("directions"),
			Notes:       c.R.FormValue("notes"),
		}

		// Tries to edit the recipe
		if _, err = c.U.Recipes().Edit(RID, newData); err == nil {
			utils.Redirect(c, "/recipes/"+strconv.Itoa(RID))
		}
	}

	return
}

// PostDeleteRecipe tries to delete a recipe
func PostDeleteRecipe(c *utils.Context) (err error) {
	var RID int
	if RID, err = getRID(c); err == nil {
		if err = c.U.Recipes().Delete(RID); err == nil {
			utils.ShowAndRedirect(c, "MSG_RECIPE_DELETED", "/recipes")
		}
	}

	return
}
