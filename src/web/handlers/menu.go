package handlers

import (
	"strconv"

	"cucinassistant/database"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func getMID(c *utils.Context) (int, error) {
	return getID(c, "MID", database.ERR_MENU_NOT_FOUND)
}

func GetMenus(c *utils.Context) (err error) {
	var menus []database.Menu

	if menus, err = c.U.Menus().GetAll(); err == nil {
		utils.RenderComponent(c, components.Menus(menus))
	}

	return
}

func GetMenusNew(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.MenusNew())
	return
}

func PostMenusNew(c *utils.Context) (err error) {
	var menu database.Menu

	if menu, err = c.U.Menus().New(c.R.FormValue("name")); err == nil {
		utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID)+"/edit")
	}

	return
}

func GetMenu(c *utils.Context) (err error) {
	var MID int
	var menu database.Menu

	if MID, err = getMID(c); err == nil {
		if menu, err = c.U.Menus().GetOne(MID); err == nil {
			utils.RenderComponent(c, components.Menu(menu))
		}
	}

	return
}

func GetMenuEdit(c *utils.Context) (err error) {
	var MID int
	var menu database.Menu

	if MID, err = getMID(c); err == nil {
		if menu, err = c.U.Menus().GetOne(MID); err == nil {
			utils.RenderComponent(c, components.MenuEdit(menu))
		}
	}

	return
}

func PostMenuEdit(c *utils.Context) (err error) {
	var MID int
	var meals [14]string

	if MID, err = getMID(c); err == nil {
		for i := 0; i < 14; i++ {
			meals[i] = c.R.FormValue("meal-" + strconv.Itoa(i))
		}

		if _, err = c.U.Menus().Replace(MID, c.R.FormValue("name"), meals); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID))
		}
	}

	return
}

func PostMenuDelete(c *utils.Context) (err error) {
	var MID int

	if MID, err = getMID(c); err == nil {
		if err = c.U.Menus().Delete(MID); err == nil {
			utils.Redirect(c, "/menus")
		}
	}

	return
}

func PostMenuDuplicate(c *utils.Context) (err error) {
	var MID int
	var menu database.Menu

	if MID, err = getMID(c); err == nil {
		if menu, err = c.U.Menus().Duplicate(MID); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(menu.MID)+"/edit")
		}
	}

	return
}
