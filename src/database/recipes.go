package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"reflect"
	"slices"

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

	// Code is a random code used by the user to share the recipe
	Code *string

	// Tags is a list of tags
	Tags []string
}

// Tag is a group of recipes that have a common tag
type Tag struct {
	// Name is the tag name
	Name string

	// Recipes are the recipes
	Recipes []Recipe
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
func (r Recipes) Edit(RID int, updated Recipe) error {
	var original Recipe
	var err error

	// Ensures the recipe (and the user) exist
	if original, err = r.GetOne(RID); err != nil {
		return err
	}

	// Ensures the stars are correct
	if updated.Stars < 0 {
		updated.Stars = 0
	} else if updated.Stars > 10 {
		updated.Stars = 10
	}

	// Checks if something has actually changed
	updated.RID = RID
	updated.Code = original.Code
	if reflect.DeepEqual(original, updated) {
		return nil
	}

	// Executes the query
	_, err = db.Exec(`UPDATE recipes SET name=$2, stars=$3, ingredients=$4, directions=$5, notes=$6 WHERE rid=$1;`,
		RID, updated.Name, updated.Stars, updated.Ingredients, updated.Directions, updated.Notes)
	if err != nil {
		if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23505" {
			return ERR_RECIPE_DUPLICATED
		} else {
			return ERR_UNKNOWN
		}
	}

	// Adds the missing tags
	for _, tag := range updated.Tags {
		if tag != "" && !slices.Contains(original.Tags, tag) {
			_, err = db.Exec(`INSERT INTO tags (name, rid) VALUES ($1, $2);`, tag, RID)
			if err != nil {
				return ERR_UNKNOWN
			}
		}
	}

	// Removes the dropped tags
	for _, tag := range original.Tags {
		if !slices.Contains(updated.Tags, tag) {
			_, err = db.Exec(`DELETE FROM tags WHERE name=$1 AND rid=$2;`, tag, RID)
			if err != nil {
				return ERR_UNKNOWN
			}
		}
	}

	return nil
}

// GetAll returns a list of recipes (with only RID and Name),
// ordered by name
func (r Recipes) GetAll() ([]Recipe, error) {
	var recipes []Recipe

	// Queries the recipes
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
	err := db.QueryRow(`SELECT rid, name, stars, ingredients, directions, notes, code FROM recipes WHERE uid=$1 AND rid=$2;`, r.uid, RID).
		Scan(&recipe.RID, &recipe.Name, &recipe.Stars, &recipe.Ingredients, &recipe.Directions, &recipe.Notes, &recipe.Code)
	if err != nil {
		return recipe, handleNoRowsError(err, r.uid, ERR_RECIPE_NOT_FOUND)
	}

	// Scans the tags
	rows, err := db.Query(`SELECT name FROM tags WHERE rid=$1 ORDER BY name;`, RID)
	if err != nil {
		return recipe, ERR_UNKNOWN
	} else {
		defer rows.Close()
		for rows.Next() {
			var tag string
			rows.Scan(&tag)
			recipe.Tags = append(recipe.Tags, tag)
		}
	}

	return recipe, nil
}

// GetPublicRecipe returns a public recipe
func GetPublicRecipe(code string) (Recipe, error) {
	var RID int
	var recipe Recipe

	// Scans the recipe
	err := db.QueryRow(`SELECT rid, name, stars, ingredients, directions, notes, code FROM recipes WHERE code=$1;`, code).
		Scan(&RID, &recipe.Name, &recipe.Stars, &recipe.Ingredients, &recipe.Directions, &recipe.Notes, &recipe.Code)
	if err != nil {
		return recipe, ERR_RECIPE_NOT_FOUND
	}

	// Scans the tags
	rows, err := db.Query(`SELECT name FROM tags WHERE rid=$1 ORDER BY name;`, RID)
	if err != nil {
		return recipe, ERR_UNKNOWN
	} else {
		defer rows.Close()
		for rows.Next() {
			var tag string
			rows.Scan(&tag)
			recipe.Tags = append(recipe.Tags, tag)
		}
	}

	return recipe, nil
}

// GetTags returns the recipes divided into tags
func (r Recipes) GetTags() ([]Tag, error) {
	var tags []Tag
	var current string

	// Queries the recipes
	var rows *sql.Rows
	rows, err := db.Query(`SELECT t.name, r.rid, r.name FROM recipes r INNER JOIN tags t ON t.rid = r.rid WHERE r.uid = $1 ORDER BY t.name, r.name;`, r.uid)
	if err != nil {
		return tags, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var recipe Recipe
		var tag string

		rows.Scan(&tag, &recipe.RID, &recipe.Name)

		if tag != current {
			tags = append(tags, Tag{Name: tag})
			current = tag
		}

		tags[len(tags)-1].Recipes = append(tags[len(tags)-1].Recipes, recipe)
	}

	// If no recipes have been found, makes sure the user exists
	if len(tags) == 0 {
		_, err := GetUser("UID", r.uid)
		return tags, err
	}

	return tags, nil
}

// NewRecipe creates a new recipe and returns its RID
func (r Recipes) New(name string) (int, error) {
	var RID int

	// Ensures the user exists
	if _, err := GetUser("UID", r.uid); err != nil {
		return RID, err
	}

	// Creates the new recipe
	err := db.QueryRow(`INSERT INTO recipes (uid, name) VALUES ($1, $2) RETURNING rid;`, r.uid, name).Scan(&RID)
	if err != nil {
		if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23505" {
			return RID, ERR_RECIPE_DUPLICATED
		} else {
			return RID, ERR_UNKNOWN
		}
	}

	return RID, nil
}

// Save creates a copy of a public recipe and returns its RID
func (r Recipes) Save(code string) (int, error) {
	var RID int

	// Gets the original
	original, err := GetPublicRecipe(code)
	if err != nil {
		return RID, err
	}

	// Creates a new one
	copiedRID, err := r.New(original.Name)
	if err != nil {
		return copiedRID, err
	}

	// Saves the content
	return copiedRID, r.Edit(copiedRID, original)
}

// Share creates a code for a recipe
func (r Recipes) Share(RID int) (string, error) {
	// Ensures the recipe (and the user) exist
	if _, err := r.GetOne(RID); err != nil {
		return "", err
	}

	for true {
		// Generates the code
		buffer := make([]byte, 4)
		rand.Read(buffer)
		code := fmt.Sprintf("%x", buffer)

		// Saves it
		_, err := db.Exec(`UPDATE recipes SET code=$2 WHERE rid=$1;`, RID, code)
		if err != nil {
			if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23505" {
				continue
			} else {
				return "", ERR_UNKNOWN
			}
		} else {
			return code, nil
		}
	}

	return "", nil
}

// Unshare deletes a recipe's code
func (r Recipes) Unshare(RID int) error {
	// Ensures the recipe (and the user) exist
	if _, err := r.GetOne(RID); err != nil {
		return err
	}

	// Saves it
	_, err := db.Exec(`UPDATE recipes SET code=NULL WHERE rid=$1;`, RID)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}
