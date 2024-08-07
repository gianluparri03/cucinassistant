package handlers

import (
	"github.com/gorilla/mux"
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetShoppingList renders /shopping_list
func GetShoppingList(c utils.Context) {
	var err error

	if user, err := database.GetUser(c.UID); err == nil {
		if list, err := user.GetShoppingList(); err == nil {
			utils.RenderPage(c, "shopping_list/view", map[string]any{"List": list})
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// GetAppendEntries renders /shopping_list/append
func GetAppendEntries(c utils.Context) {
	utils.RenderPage(c, "shopping_list/append", nil)
}

// PostAppendEntries tries to append the given entries
// to the shopping list
func PostAppendEntries(c utils.Context) {
	var names []string

	// Insert all the values whose key is entry-X-name
	c.R.ParseForm()
	for key, values := range c.R.PostForm {
		if strings.HasPrefix(key, "entry-") && strings.HasSuffix(key, "-name") {
			if len(values) > 0 && values[0] != "" {
				names = append(names, values[0])
			}
		}
	}

	// Tries to append them to the list
	var err error
	if user, err := database.GetUser(c.UID); err == nil {
		if err = user.AppendEntries(names...); err == nil {
			utils.Redirect(c, "/shopping_list")
			return
		}
	}

	// Handles errors
	utils.ShowMessage(c, err.Error(), "")
}

// PostToggleEntry tries to toggle an entry in the shopping list
func PostToggleEntry(c utils.Context) {
	// Retrieves the EID
	EID, err := strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_ENTRY_NOT_FOUND.Error(), "/shopping_list")
		return
	}

	// Tries to toggle the entry
	var user database.User
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.ToggleEntry(EID); err == nil {
			utils.Redirect(c, "/shopping_list")
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "/shopping_list")
}

// PostClearEntries tries to deletes all the marked entries
func PostClearEntries(c utils.Context) {
	// Tries to clear the list
	var user database.User
	var err error
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.ClearEntries(); err == nil {
			utils.ShowMessage(c, "Lista svuotata con successo", "/shopping_list")
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "/shopping_list")
}

// GetEditEntry renders /shopping_list/{EID}/edit
func GetEditEntry(c utils.Context) {
	// Retrieves the EID
	EID, err := strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_ENTRY_NOT_FOUND.Error(), "/shopping_list")
		return
	}

	// Retrieves the entry's name and renders the page
	if user, err := database.GetUser(c.UID); err == nil {
		if entry, err := user.GetEntry(EID); err == nil {
			utils.RenderPage(c, "shopping_list/edit", map[string]any{"Name": entry.Name})
			return
		}
	}

	utils.ShowMessage(c, err.Error(), "")
}

// PostEditEntry tries to change an entry's name
func PostEditEntry(c utils.Context) {
	var err error

	// Retrieves the EID
	EID, err := strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		utils.ShowMessage(c, database.ERR_ENTRY_NOT_FOUND.Error(), "/shopping_list")
		return
	}

	// Retrieves the new name
	newName := c.R.FormValue("name")

	// Tries to toggle the entry
	if user, err := database.GetUser(c.UID); err == nil {
		if err = user.EditEntry(EID, newName); err == nil {
			utils.Redirect(c, "/shopping_list")
			return
		}
	}

	// Handles the error
	utils.ShowMessage(c, err.Error(), "")
}
