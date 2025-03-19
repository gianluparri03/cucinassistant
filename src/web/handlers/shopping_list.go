package handlers

import (
	"strings"

	"cucinassistant/database"
	"cucinassistant/langs"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

// getEID returns the EID written in the url
func getEID(c *utils.Context) (int, error) {
	return getID(c, "EID", database.ERR_ENTRY_NOT_FOUND)
}

// GetShoppingList renders /shopping_list
func GetShoppingList(c *utils.Context) (err error) {
	var list map[int]database.Entry

	if list, err = c.U.ShoppingList().GetAll(); err == nil {
		utils.RenderComponent(c, components.ShoppingListView(list))
	}

	return
}

// GetAppendEntries renders /shopping_list/append
func GetAppendEntries(c *utils.Context) error {
	utils.RenderComponent(c, components.ShoppingListAppend())
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
	if err = c.U.ShoppingList().Append(names...); err == nil {
		utils.Redirect(c, "/shopping_list")
	}

	return
}

// PostToggleEntry tries to toggle an entry in the shopping list
func PostToggleEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = getEID(c); err == nil {
		// Tries to toggle the entry
		if err = c.U.ShoppingList().Toggle(EID); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}

// PostClearShoppingList tries to deletes all the marked entries
func PostClearShoppingList(c *utils.Context) (err error) {
	// Tries to clear the list
	if err = c.U.ShoppingList().Clear(); err == nil {
		utils.ShowMessage(c, langs.STR_SHOPPINGLIST_EMPTIED, "/shopping_list")
	}

	return
}

// GetEditEntry renders /shopping_list/{EID}/edit
func GetEditEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = getEID(c); err == nil {
		// Retrieves the entry's name and renders the page
		var entry database.Entry
		if entry, err = c.U.ShoppingList().GetOne(EID); err == nil {
			utils.RenderComponent(c, components.ShoppingListEdit(entry.Name))
		}
	}

	return
}

// PostEditEntry tries to change an entry's name
func PostEditEntry(c *utils.Context) (err error) {
	// Retrieves the EID
	var EID int
	if EID, err = getEID(c); err == nil {
		// Retrieves the new name
		newName := c.R.FormValue("name")

		// Tries to toggle the entry
		if err = c.U.ShoppingList().Edit(EID, newName); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}
