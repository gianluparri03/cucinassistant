package handlers

import (
	"strconv"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// getMID returns the MID written in the url
func getMID(c *utils.Context) (MID int, err error) {
	return getID(c, "MID", database.ERR_MENU_NOT_FOUND)
}

// GetMenus renders /menus
func GetMenus(c *utils.Context) (err error) {
	var menus []database.Menu
	if menus, err = c.U.GetMenus(); err == nil {
		utils.RenderPage(c, "menu/dashboard", map[string]any{"Menus": menus})
	}

	return
}

// PostNewMenu tries to create a new menu
func PostNewMenu(c *utils.Context) (err error) {
	var menu database.Menu
	if menu, err = c.U.NewMenu(); err == nil {
		utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID)+"/edit")
	}

	return
}

// GetMenu renders /menus/{MID}
func GetMenu(c *utils.Context) (err error) {
	// Retrieves the MID
	var MID int
	if MID, err = getMID(c); err == nil {
		// Retrieves the menu
		var menu database.Menu
		if menu, err = c.U.GetMenu(MID); err == nil {
			utils.RenderPage(c, "menu/view", map[string]any{"Menu": menu})
		}
	}

	return
}

// GetEditMenu renders /menus/{MID}/edit
func GetEditMenu(c *utils.Context) (err error) {
	// Retrieves the MID
	var MID int
	if MID, err = getMID(c); err == nil {
		// Retrieves the menu
		var menu database.Menu
		if menu, err = c.U.GetMenu(MID); err == nil {
			utils.RenderPage(c, "menu/edit", map[string]any{"Menu": menu})
		}
	}

	return
}

// PostEditMenu tries to replace a menu
func PostEditMenu(c *utils.Context) (err error) {
	// Retrieves the MID
	var MID int
	if MID, err = getMID(c); err == nil {
		// Retrieves the new data
		var meals [14]string
		for i := 0; i < 14; i++ {
			meals[i] = c.R.FormValue("meal-" + strconv.Itoa(i))
		}

		// Tries to replace the menu
		if _, err = c.U.ReplaceMenu(MID, c.R.FormValue("name"), meals); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID))
		}
	}

	return
}

// PostDuplicateMenu tries to duplicate a menu
func PostDuplicateMenu(c *utils.Context) (err error) {
	// Retrieves the MID
	var MID int
	if MID, err = getMID(c); err == nil {
		// Tries to duplicate the menu
		var menu database.Menu
		if menu, err = c.U.DuplicateMenu(MID); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID)+"/edit")
		}
	}

	return
}

// PostDeleteMenu tries to delete a menu
func PostDeleteMenu(c *utils.Context) (err error) {
	// Retrieves the MID
	var MID int
	if MID, err = getMID(c); err == nil {
		// Tries to delete the menu
		if err = c.U.DeleteMenu(MID); err == nil {
			utils.ShowAndRedirect(c, "Menù eliminato con successo", "/menus")
		}
	}

	return
}
