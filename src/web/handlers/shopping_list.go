package handlers

import (
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func getEID(c *utils.Context) (int, error) {
	return getID(c, "EID", database.ERR_ENTRY_NOT_FOUND)
}

func GetShoppingList(c *utils.Context) (err error) {
	var list map[int]database.Entry

	if list, err = c.U.ShoppingList().GetAll(); err == nil {
		utils.RenderComponent(c, components.ShoppingList(list))
	}

	return
}

func GetShoppingListAppend(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.ShoppingListAppend())
	return
}

func PostShoppingListAppend(c *utils.Context) (err error) {
	var names []string
	c.R.ParseForm()

	for key, values := range c.R.PostForm {
		if strings.HasPrefix(key, "entry-") && strings.HasSuffix(key, "-name") {
			if len(values) > 0 && values[0] != "" {
				names = append(names, values[0])
			}
		}
	}

	if err = c.U.ShoppingList().Append(names...); err == nil {
		utils.Redirect(c, "/shopping_list")
	}

	return
}

func PostShoppingListClear(c *utils.Context) (err error) {
	if err = c.U.ShoppingList().Clear(); err == nil {
		utils.Redirect(c, "/shopping_list")
	}

	return
}

func GetEntryEdit(c *utils.Context) (err error) {
	var EID int
	var entry database.Entry

	if EID, err = getEID(c); err == nil {
		if entry, err = c.U.ShoppingList().GetOne(EID); err == nil {
			utils.RenderComponent(c, components.EntryEdit(entry.Name))
		}
	}

	return
}

func PostEntryEdit(c *utils.Context) (err error) {
	var EID int

	if EID, err = getEID(c); err == nil {
		newName := c.R.FormValue("name")

		if err = c.U.ShoppingList().Edit(EID, newName); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}

func PostEntryToggle(c *utils.Context) (err error) {
	var EID int

	if EID, err = getEID(c); err == nil {
		if err = c.U.ShoppingList().Toggle(EID); err == nil {
			utils.Redirect(c, "/shopping_list")
		}
	}

	return
}
