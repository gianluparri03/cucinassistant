package database

import (
	"strings"
)

const (
	// menuDefaultName is the name given to new menus
	menuDefaultName = "Menù"

	// mealSeparator is used to separate meals when packed
	mealSeparator = ";"

	// duplicatedMenuSuffix is added at the end of the name when
	// duplicating a menu
	duplicatedMenuSuffix = " (copia)"
)

// packMeals packs the 14 meals in a string
func packMeals(meals [14]string) string {
	var sb strings.Builder

	for index, meal := range meals {
		sb.WriteString(meal)

		if index < 13 {
			sb.WriteString(mealSeparator)
		}
	}

	return sb.String()
}

// unpackMeals unpacks a string in an array of meals
func unpackMeals(meals string) (out [14]string) {
	for index, meal := range strings.Split(meals, mealSeparator) {
		if index == 14 {
			break
		}

		out[index] = meal
	}

	return
}
