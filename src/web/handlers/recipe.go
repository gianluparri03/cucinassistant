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

func getRID(c *utils.Context) (int, error) {
	return getID(c, "RID", database.ERR_RECIPE_NOT_FOUND)
}

func GetPublicRecipe(c *utils.Context) (err error) {
	var recipe database.Recipe

	code := mux.Vars(c.R)["code"]
	if recipe, err = database.GetPublicRecipe(code); err == nil {
		utils.RenderComponent(c, components.Recipe(recipe, configs.BaseURL))
	}

	return
}

func PostPublicRecipeSave(c *utils.Context) (err error) {
	var recipe database.Recipe

	if recipe, err = c.U.Recipes().Save(mux.Vars(c.R)["code"]); err == nil {
		utils.ShowMessage(c, langs.STR_RECIPE_COPIED, "/recipes/"+strconv.Itoa(recipe.RID))
	}

	return
}

func GetRecipes(c *utils.Context) (err error) {
	var recipes []database.Recipe

	if recipes, err = c.U.Recipes().GetAll(); err == nil {
		utils.RenderComponent(c, components.Recipes(recipes))
	}

	return
}

func GetRecipesNew(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.RecipesNew())
	return
}

func PostRecipesNew(c *utils.Context) (err error) {
	var recipe database.Recipe

	if recipe, err = c.U.Recipes().New(c.R.FormValue("name")); err == nil {
		utils.Redirect(c, "/recipes/"+strconv.Itoa(recipe.RID)+"/edit")
	}

	return
}

func GetRecipe(c *utils.Context) (err error) {
	var RID int
	var recipe database.Recipe

	if RID, err = getRID(c); err == nil {
		if recipe, err = c.U.Recipes().GetOne(RID); err == nil {
			utils.RenderComponent(c, components.Recipe(recipe, configs.BaseURL))
		}
	}

	return
}

func GetRecipeEdit(c *utils.Context) (err error) {
	var RID int
	var recipe database.Recipe

	if RID, err = getRID(c); err == nil {
		if recipe, err = c.U.Recipes().GetOne(RID); err == nil {
			utils.RenderComponent(c, components.RecipeEdit(recipe))
		}
	}

	return
}

func PostRecipeEdit(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		stars, _ := strconv.ParseFloat(c.R.FormValue("stars"), 32)
		newData := database.Recipe{
			Name:        c.R.FormValue("name"),
			Stars:       int(stars * 2),
			Ingredients: c.R.FormValue("ingredients"),
			Directions:  c.R.FormValue("directions"),
			Notes:       c.R.FormValue("notes"),
		}

		if _, err = c.U.Recipes().Edit(RID, newData); err == nil {
			utils.Redirect(c, "/recipes/"+strconv.Itoa(RID))
		}
	}

	return
}

func PostRecipeDelete(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		if err = c.U.Recipes().Delete(RID); err == nil {
			utils.ShowMessage(c, langs.STR_RECIPE_DELETED, "/recipes")
		}
	}

	return
}

func GetRecipeShare(c *utils.Context) (err error) {
	var RID int
	var recipe database.Recipe

	if RID, err = getRID(c); err == nil {
		if recipe, err = c.U.Recipes().GetOne(RID); err == nil {
			utils.RenderComponent(c, components.RecipeShare(recipe, configs.BaseURL))
		}
	}

	return
}

func PostRecipeShare(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		if _, err = c.U.Recipes().Share(RID); err == nil {
			utils.Redirect(c, "/recipes/"+strconv.Itoa(RID)+"/share")
		}
	}

	return
}

func PostRecipeUnshare(c *utils.Context) (err error) {
	var RID int

	if RID, err = getRID(c); err == nil {
		if err = c.U.Recipes().Unshare(RID); err == nil {
			utils.Redirect(c, "/recipes/"+strconv.Itoa(RID)+"/share")
		}
	}

	return
}
