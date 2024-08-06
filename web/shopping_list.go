package web

import (
	"github.com/gorilla/mux"
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// appendEntries tries to append the given entries
// to the shopping list
func appendEntries(c utils.Context) {
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
	var user database.User
	var err error
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.AppendEntries(names...); err == nil {
			utils.ShowError(c, "Elementi aggiunti con successo", "/shopping_list")
			return
		}
	}

	// Handles errors
	utils.ShowError(c, err.Error(), "")
}

// toggleEntry tries to toggle an entry in the shopping list
func toggleEntry(c utils.Context) {
	// Retrieves the EID
	EID, err := strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		utils.ShowError(c, database.ERR_ENTRY_NOT_FOUND.Error(), "/shopping_list")
		return
	}

	// Tries to toggle the entry
	var user database.User
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.ToggleEntry(EID); err == nil {
			return
		}
	}

	// Handles the error
	utils.ShowError(c, err.Error(), "/shopping_list")
}

// clearEntries tries to deletes all the marked entries
func clearEntries(c utils.Context) {
	// Tries to clear the list
	var user database.User
	var err error
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.ClearEntries(); err == nil {
			utils.ShowError(c, "Lista svuotata con successo", "/shopping_list")
			return
		}
	}

	// Handles the error
	utils.ShowError(c, err.Error(), "/shopping_list")
}

// editEntry tries to change an entry's name
func editEntry(c utils.Context) {
	// Retrieves the EID
	EID, err := strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		utils.ShowError(c, database.ERR_ENTRY_NOT_FOUND.Error(), "/shopping_list")
		return
	}

	// Retrieves the new name
	newName := c.R.FormValue("name")

	// Tries to toggle the entry
	var user database.User
	if user, err = database.GetUser(c.UID); err == nil {
		if err = user.EditEntry(EID, newName); err == nil {
			utils.ShowError(c, "Elemento modificato con successo", "/shopping_list")
			return
		}
	}

	// Handles the error
	utils.ShowError(c, err.Error(), "")
}
