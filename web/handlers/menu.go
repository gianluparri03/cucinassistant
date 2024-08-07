package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetMenus renders /menus
func GetMenus(c utils.Context) {
	var err error
	if user, err := database.GetUser(c.UID); err == nil {
		if menus, err := user.GetMenus(); err == nil {
			utils.RenderPage(c, "menu/dashboard", map[string]any{"Menus": menus})
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// PostNewMenu tries to create a new menu
func PostNewMenu(c utils.Context) {
	var err error

	if user, err := database.GetUser(c.UID); err == nil {
		if menu, err := user.NewMenu(); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID))
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// GetMenu renders /menus/{MID}
func GetMenu(c utils.Context) {
	// Retrieves the MID
	MID, err := strconv.Atoi(mux.Vars(c.R)["MID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_MENU_NOT_FOUND.Error(), "/menus")
		return
	}

	// Retrieves the menu
	if user, err := database.GetUser(c.UID); err == nil {
		if menu, err := user.GetMenu(MID); err == nil {
			utils.RenderPage(c, "menu/view", map[string]any{"Menu": menu})
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// GetEditMenu renders /menus/{MID}/edit
func GetEditMenu(c utils.Context) {
	// Retrieves the MID
	MID, err := strconv.Atoi(mux.Vars(c.R)["MID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_MENU_NOT_FOUND.Error(), "/menus")
		return
	}

	// Retrieves the menu
	if user, err := database.GetUser(c.UID); err == nil {
		if menu, err := user.GetMenu(MID); err == nil {
			utils.RenderPage(c, "menu/edit", map[string]any{"Menu": menu})
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// PostEditMenu tries to replace a menu
func PostEditMenu(c utils.Context) {
	var err error

	// Retrieves the MID
	MID, err := strconv.Atoi(mux.Vars(c.R)["MID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_MENU_NOT_FOUND.Error(), "/menus")
		return
	}

	// Retrieves the new data
	var data database.Menu
	data.Name = c.R.FormValue("name")
	for i := 0; i < 14; i++ {
		data.Meals[i] = c.R.FormValue("meal-" + strconv.Itoa(i))
	}

	// Tries to toggle the entry
	if user, err := database.GetUser(c.UID); err == nil {
		if err = user.ReplaceMenu(MID, data); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID))
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "")
}

// PostDuplicateMenu tries to duplicate a menu
func PostDuplicateMenu(c utils.Context) {
	var err error

	// Retrieves the MID
	MID, err := strconv.Atoi(mux.Vars(c.R)["MID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_MENU_NOT_FOUND.Error(), "/menus")
		return
	}

	// Tries to duplicate it
	if user, err := database.GetUser(c.UID); err == nil {
		if menu, err := user.DuplicateMenu(MID); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID)+"/edit")
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "")
}

// PostDeleteMenu tries to delete a menu
func PostDeleteMenu(c utils.Context) {
	// Retrieves the MID
	MID, err := strconv.Atoi(mux.Vars(c.R)["MID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_MENU_NOT_FOUND.Error(), "/menus")
		return
	}

	// Tries to duplicate it
	if user, err := database.GetUser(c.UID); err == nil {
		if err = user.DeleteMenu(MID); err == nil {
			utils.ShowMessage(c, "Menù eliminato con successo", "/menus")
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "")
}
