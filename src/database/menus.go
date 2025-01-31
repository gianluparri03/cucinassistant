package database

import (
	"database/sql"
	"log/slog"
	"strings"
)

const (
	// menuDefaultName is the name given to new menus
	menuDefaultName = "Nuovo Menù"

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

// Menus is used to manage all the menus
type Menus struct {
	uid int
}

// Menus returns the menus manager for the user
func (u User) Menus() Menus {
	return Menus{uid: u.UID}
}

// Menu is a collection of 14 meals, from monday lunch (0) to sunday dinner (13)
type Menu struct {
	// MID is the Menu ID
	MID int

	// Name is the name of the menu
	Name string

	// Meals contains all the 14 meals
	Meals [14]string
}

// GetAll returns a list of menus (meals not included)
func (m Menus) GetAll() ([]Menu, error) {
	var menus []Menu

	// Queries the entries
	var rows *sql.Rows
	rows, err := db.Query(`SELECT mid, name FROM menus WHERE uid=$1;`, m.uid)
	if err != nil {
		slog.Error("while retrieving menus:", "err", err)
		return menus, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var m Menu
		rows.Scan(&m.MID, &m.Name)
		menus = append(menus, m)
	}

	// If no menus have been found, makes sure the user exists
	if len(menus) == 0 {
		_, err := GetUser("UID", m.uid)
		return menus, err
	}

	return menus, nil
}

// GetOne returns a specific menu, with the meals
func (m Menus) GetOne(MID int) (Menu, error) {
	var menu Menu
	var meals string

	// Scans the menu
	err := db.QueryRow(`SELECT mid, name, meals FROM menus WHERE uid=$1 AND mid=$2;`, m.uid, MID).Scan(&menu.MID, &menu.Name, &meals)
	if err != nil {
		return menu, handleNoRowsError(err, m.uid, ERR_MENU_NOT_FOUND, "retrieving menu")
	}

	// Unpacks the meals
	menu.Meals = unpackMeals(meals)
	return menu, nil
}

// New creates a new menu
func (m Menus) New() (Menu, error) {
	// Ensures the user exists
	if _, err := GetUser("UID", m.uid); err != nil {
		return Menu{}, err
	}

	// Adds the new menu
	menu := Menu{Name: menuDefaultName}
	err := db.QueryRow(`INSERT INTO menus (uid, name, meals) VALUES ($1, $2, $3) RETURNING mid;`, m.uid, menu.Name, packMeals(menu.Meals)).Scan(&menu.MID)
	if err != nil {
		slog.Error("while creating new menu:", "err", err)
		return Menu{}, ERR_UNKNOWN
	}

	return menu, nil
}

// Replace replaces the menu's name and all of its meals
func (m Menus) Replace(MID int, newName string, newMeals [14]string) (Menu, error) {
	// Ensures the menu (and the user) exist
	if _, err := m.GetOne(MID); err != nil {
		return Menu{}, err
	}

	// Executes the query
	_, err := db.Exec(`UPDATE menus SET name=$3, meals=$4 WHERE uid=$1 AND mid=$2;`, m.uid, MID, newName, packMeals(newMeals))
	if err != nil {
		slog.Error("while replacing menu:", "err", err)
		err = ERR_UNKNOWN
		return Menu{}, err
	}

	// Returns the new menu
	return m.GetOne(MID)
}

// Delete deletes a menu
func (m Menus) Delete(MID int) error {
	// Deletes the menu
	res, err := db.Exec(`DELETE FROM menus WHERE uid=$1 AND mid=$2;`, m.uid, MID)
	if err != nil {
		slog.Error("while deleting menu:", "err", err)
		return ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the menu (and the user) exist
		_, err := m.GetOne(MID)
		return err
	}

	return nil
}

// Duplicate creates a copy of a menu
func (m Menus) Duplicate(MID int) (Menu, error) {
	// Gets the source menu
	if src, err := m.GetOne(MID); err != nil {
		return Menu{}, err
		// Creates a new one
	} else if menu, err := m.New(); err != nil {
		return Menu{}, err
		// Copies data from the srcMenu
	} else {
		return m.Replace(menu.MID, src.Name+duplicatedMenuSuffix, src.Meals)
	}
}
