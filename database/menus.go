package database

import (
	"database/sql"
	"log/slog"
	"strings"
)

// Menu is a collection of 14 meals, from monday lunch (0) to sunday dinner (13)
type Menu struct {
	// MID is the Menu ID
	MID int

	// Name is the name of the menu
	Name string

	// Meals contains all the 14 meals
	Meals [14]string
}

// GetMenus returns a list of menus (meals not included)
func (u *User) GetMenus() (menus []*Menu, err error) {
	menus = []*Menu{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = db.Query(`SELECT mid, name FROM menus WHERE uid=$1;`, u.UID)
	if err != nil {
		slog.Error("while retrieving menus:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var m Menu
		rows.Scan(&m.MID, &m.Name)
		menus = append(menus, &m)
	}

	// If no menus have been found, makes sure the user exists
	if len(menus) == 0 {
		if _, err = GetUser("UID", u.UID); err != nil {
			menus = nil
		}
	}

	return
}

// GetMenu returns a specific menu, with the meals
func (u *User) GetMenu(MID int) (menu *Menu, err error) {
	menu = &Menu{}
	var meals string

	// Scans the menu
	err = db.QueryRow(`SELECT mid, name, meals FROM menus WHERE uid=$1 AND mid=$2;`, u.UID, MID).Scan(&menu.MID, &menu.Name, &meals)
	if err != nil {
		menu = nil

		if strings.HasSuffix(err.Error(), "no rows in result set") {
			// Makes sure the user exists if the menu has not been found
			if _, err = GetUser("UID", u.UID); err == nil {
				err = ERR_MENU_NOT_FOUND
			}
		} else {
			slog.Error("while retrieving menu:", "err", err)
			err = ERR_UNKNOWN
		}

		return
	}

	// Unpacks the meals
	menu.Meals = unpackMeals(meals)
	return
}

// NewMenu creates a new menu
func (u *User) NewMenu() (menu *Menu, err error) {
	// Ensures the user exists
	_, err = GetUser("UID", u.UID)
	if err != nil {
		return
	}

	// Uses the default name
	menu = &Menu{}
	menu.Name = menuDefaultName

	// Adds the new menu
	_, err = db.Exec(`INSERT INTO menus (uid, name, meals) VALUES ($1, $2, $3);`, u.UID, menu.Name, packMeals(menu.Meals))
	if err != nil {
		slog.Error("while creating new menu:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Fetches the MID
	db.QueryRow(`SELECT MAX(mid) FROM menus WHERE uid=$1;`, u.UID).Scan(&menu.MID)
	return
}

// ReplaceMenu replaces the menu's name and all of its meals
func (u *User) ReplaceMenu(MID int, newName string, newMeals [14]string) (menu *Menu, err error) {
	// Ensures the menu (and the user) exist
	if _, err = u.GetMenu(MID); err != nil {
		return
	}

	// Executes the query
	_, err = db.Exec(`INSERT INTO menus (uid, mid, name, meals) VALUES ($1, $2, $3, $4)
                      ON CONFLICT (mid) DO UPDATE
                      SET name=$3, meals=$4;`, u.UID, MID, newName, packMeals(newMeals))
	if err != nil {
		slog.Error("while replacing menu:", "err", err)
		err = ERR_UNKNOWN
	}

	// Returns the new menu
	menu, err = u.GetMenu(MID)
	return
}

// DeleteMenu deletes a menu
func (u *User) DeleteMenu(MID int) (err error) {
	// Executes the query
	res, err := db.Exec(`DELETE FROM menus WHERE uid=$1 AND mid=$2;`, u.UID, MID)
	if err != nil {
		slog.Error("while deleting menu:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the menu (and the user) exist
		if _, err = GetUser("UID", u.UID); err == nil {
			err = ERR_MENU_NOT_FOUND
		}
	}

	return
}

// DuplicateMenu creates a copy of a menu
func (u *User) DuplicateMenu(MID int) (menu *Menu, err error) {
	var src *Menu

	// Gets the source menu
	if src, err = u.GetMenu(MID); err != nil {
		return
	}

	// Creates a new one
	if menu, err = u.NewMenu(); err != nil {
		return
	}

	// Copies data from the srcMenu
	menu, err = u.ReplaceMenu(menu.MID, src.Name+duplicatedMenuSuffix, src.Meals)
	return
}
