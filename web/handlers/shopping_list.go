package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetEID returns the EID written in the url
func GetEID(c *utils.Context) (EID int, err error) {
	// Fetches the EID from the url, then converts
	// it to an int
	EID, err = strconv.Atoi(mux.Vars(c.R)["EID"])
	if err != nil {
		err = database.ERR_ENTRY_NOT_FOUND
	}

	return
}

// GetShoppingList renders /shopping_list
func GetShoppingList(c *utils.Context) (err error) {
	var list map[int]*database.Entry
	if list, err = c.U.GetShoppingList(); err == nil {
		utils.RenderPage(c, "shopping_list/view", map[string]any{"List": list})
	}

	return
}

// GetAppendEntries renders /shopping_list/append
func GetAppendEntries(c *utils.Context) error {
	utils.RenderPage(c, "shopping_list/append", nil)
	return nil
}

// PostAppendEntries tries to append the given entries
// to the shopping list
func PostAppendEntries(c *utils.Context) (err error) {
	var names []string
	c.R.ParseForm()

	// Insert all the values whose key is entry-X-name
	for key, values := range c.R.PostForm {
		if strings.HasPrefix(key, "entry-") && strings.HasSuffix(key, "-name") {
			if len(values) > 0 && values[0] != "" {
				names = append(names, values[0])
			}
		}
	}

	// Tries to append them to the list
	if err = c.U.AppendEntries(names...); err == nil {
		utils.Redirect(c, "/shopping_list")
	}

	return
}

// PostToggleEntry tries to toggle an entry in the shopping list
func PostToggleEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = GetEID(c); err == nil {
		// Tries to toggle the entry
		if err = c.U.ToggleEntry(EID); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}

// PostClearEntries tries to deletes all the marked entries
func PostClearEntries(c *utils.Context) (err error) {
	// Tries to clear the list
	if err = c.U.ClearEntries(); err == nil {
		utils.ShowAndRedirect(c, "Lista svuotata con successo", "/shopping_list")
	}

	return
}

// GetEditEntry renders /shopping_list/{EID}/edit
func GetEditEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = GetEID(c); err == nil {
		// Retrieves the entry's name and renders the page
		var entry *database.Entry
		if entry, err = c.U.GetEntry(EID); err == nil {
			utils.RenderPage(c, "shopping_list/edit", map[string]any{"Name": entry.Name})
		}
	}

	fmt.Println(err)

	return
}

// PostEditEntry tries to change an entry's name
func PostEditEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = GetEID(c); err == nil {
		// Retrieves the new name
		newName := c.R.FormValue("name")

		// Tries to toggle the entry
		if err = c.U.EditEntry(EID, newName); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}
