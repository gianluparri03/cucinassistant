package database

import (
	"database/sql"
	"reflect"

	"github.com/lib/pq"
)

// Recipe contains a RID, a name, some ingredients,
// some directives, some notes and a number of stars
type Recipe struct {
	// RID is the Recipe ID
	RID int

	// Name is the name of the recipe
	Name string

	// Stars is the number of stars the recipe
	// has (0 <= Stars <= 5)
	Stars int

	// Ingredients is a text containing the ingredients
	Ingredients string

	// Directions is a text containing the directions
	Directions string

	// Notes are additional notes to the recipe
	Notes string
}

// Recipes is used to manage all the recipes
type Recipes struct {
	uid int
}

// Recipes returns the recipe manager for the user
func (u User) Recipes() Recipes {
	return Recipes{uid: u.UID}
}

// Delete deletes a recipe
func (r Recipes) Delete(RID int) error {
	res, err := db.Exec(`DELETE FROM recipes WHERE uid=$1 AND rid=$2;`, r.uid, RID)
	if err != nil {
		return ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the recipe (and the user) exist
		_, err := r.GetOne(RID)
		return err
	}

	return nil
}

// Edit replaces all the recipes's data, except for the RID
func (r Recipes) Edit(RID int, updated Recipe) (Recipe, error) {
	var original Recipe
	var err error

	// Ensures the recipe (and the user) exist
	if original, err = r.GetOne(RID); err != nil {
		return Recipe{}, err
	}

	// Ensures the stars are correct
	if updated.Stars < 0 {
		updated.Stars = 0
	} else if updated.Stars > 10 {
		updated.Stars = 10
	}

	// Checks if something has actually changed
	updated.RID = RID
	if reflect.DeepEqual(original, updated) {
		return original, nil
	}

	// Executes the query
	_, err = db.Exec(`UPDATE recipes SET name=$3, stars=$4, ingredients=$5, directions=$6, notes=$7 WHERE uid=$1 AND rid=$2;`,
		r.uid, RID, updated.Name, updated.Stars, updated.Ingredients, updated.Directions, updated.Notes)
	if err != nil {
		if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23505" {
			return Recipe{}, ERR_RECIPE_DUPLICATED
		} else {
			return Recipe{}, ERR_UNKNOWN
		}
	}

	return updated, nil
}

// GetAll returns a list of recipes (with only RID and Name),
// ordered by name
func (r Recipes) GetAll() ([]Recipe, error) {
	var recipes []Recipe

	// Queries the recipies
	var rows *sql.Rows
	rows, err := db.Query(`SELECT rid, name FROM recipes WHERE uid=$1 ORDER BY name;`, r.uid)
	if err != nil {
		return recipes, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var r Recipe
		rows.Scan(&r.RID, &r.Name)
		recipes = append(recipes, r)
	}

	// If no recipes have been found, makes sure the user exists
	if len(recipes) == 0 {
		_, err := GetUser("UID", r.uid)
		return recipes, err
	}

	return recipes, nil
}

// GetOne returns a specific recipe
func (r Recipes) GetOne(RID int) (Recipe, error) {
	var recipe Recipe

	// Scans the recipe
	err := db.QueryRow(`SELECT rid, name, stars, ingredients, directions, notes FROM recipes WHERE uid=$1 AND rid=$2;`, r.uid, RID).
		Scan(&recipe.RID, &recipe.Name, &recipe.Stars, &recipe.Ingredients, &recipe.Directions, &recipe.Notes)
	if err != nil {
		return recipe, handleNoRowsError(err, r.uid, ERR_RECIPE_NOT_FOUND)
	}

	return recipe, nil
}

// NewRecipe creates a new recipe
func (r Recipes) New(name string) (Recipe, error) {
	// Ensures the user exists
	if _, err := GetUser("UID", r.uid); err != nil {
		return Recipe{}, err
	}

	// Creates the new recipe
	recipe := Recipe{Name: name}
	err := db.QueryRow(`INSERT INTO recipes (uid, name) VALUES ($1, $2) RETURNING rid;`, r.uid, recipe.Name).Scan(&recipe.RID)
	if err != nil {
		if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23505" {
			return Recipe{}, ERR_RECIPE_DUPLICATED
		} else {
			return Recipe{}, ERR_UNKNOWN
		}
	}

	return recipe, nil
}
