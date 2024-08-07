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
func (u *User) GetMenus() (menus []Menu, err error) {
	menus = []Menu{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = DB.Query(`SELECT mid, name FROM menus WHERE uid=?;`, u.UID)
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
		menus = append(menus, m)
	}

	// If no menus have been found, makes sure the user exists
	if len(menus) == 0 {
		_, err = GetUser(u.UID)
	}

	return
}

// GetMenu returns a specific menu, with the meals
func (u *User) GetMenu(MID int) (menu Menu, err error) {
	var meals string

	// Scans the menu
	err = DB.QueryRow(`SELECT mid, name, meals FROM menus WHERE uid=? AND mid=?;`, u.UID, MID).Scan(&menu.MID, &menu.Name, &meals)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			// Makes sure the user exists if the menu has not been found
			if _, err = GetUser(u.UID); err == nil {
				err = ERR_MENU_NOT_FOUND
			}
		} else {
			slog.Error("while retrieving menu:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	// Unpacks the meals
	menu.Meals = unpackMeals(meals)

	return
}

// NewMenu creates a new menu
func (u *User) NewMenu() (menu Menu, err error) {
	// Ensures the user exists
	_, err = GetUser(u.UID)
	if err != nil {
		return
	}

	// Uses the default name
	menu.Name = menuDefaultName

	// Adds the new menu
	_, err = DB.Exec(`INSERT INTO menus (uid, mid, name, meals)
	                  SELECT ?, IFNULL(MAX(mid), 0) + 1, ?, ? FROM menus WHERE uid=?;`, u.UID, menu.Name, packMeals(menu.Meals), u.UID)
	if err != nil {
		slog.Error("while creating new menu:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Fetches the MID
	DB.QueryRow(`SELECT MAX(mid) FROM menus WHERE uid=?;`, u.UID).Scan(&menu.MID)
	return
}

// ReplaceMenu replaces the menu's name and all of its meals
func (u *User) ReplaceMenu(MID int, newData Menu) (err error) {
	// Ensures the menu (and the user) exist
	if _, err = u.GetMenu(MID); err != nil {
		return
	}

	// Prepares the values for the query
	values := []any{u.UID, MID, newData.Name}
	for i := 0; i < 14; i++ {
		values = append(values, newData.Meals[i])
	}

	// Executes the query
	_, err = DB.Exec(`REPLACE INTO menus (uid, mid, name, meals) VALUES (?, ?, ?, ?);`, u.UID, MID, newData.Name, packMeals(newData.Meals))
	if err != nil {
		slog.Error("while replacing menu:", "err", err)
		err = ERR_UNKNOWN
	}

	return
}

// DeleteMenu deletes a menu
func (u *User) DeleteMenu(MID int) (err error) {
	// Executes the query
	res, err := DB.Exec(`DELETE FROM menus WHERE uid=? AND mid=?;`, u.UID, MID)
	if err != nil {
		slog.Error("while deleting menu:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the menu (and the user) exist
		if _, err = GetUser(u.UID); err == nil {
			err = ERR_MENU_NOT_FOUND
		}
	}

	return
}

// DuplicateMenu creates a copy of a menu
func (u *User) DuplicateMenu(MID int) (dstMenu Menu, err error) {
	var srcMenu Menu

	// Gets the source menu
	if srcMenu, err = u.GetMenu(MID); err != nil {
		return
	}

	// Creates a new one
	if dstMenu, err = u.NewMenu(); err != nil {
		return
	}

	// Copies the content
	dstMenu.Name = srcMenu.Name
	dstMenu.Meals = srcMenu.Meals

	// Updates it in the database
	err = u.ReplaceMenu(dstMenu.MID, dstMenu)
	return
}
