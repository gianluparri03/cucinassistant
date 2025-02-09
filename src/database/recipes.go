package database

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
