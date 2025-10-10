package handlers

import (
	"strconv"
	"slices"
	"strings"

	"cucinassistant/database"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func getMID(c *utils.Context) (int, error) {
	return getID(c, "MID", database.ERR_MENU_NOT_FOUND)
}

func getDPos(c *utils.Context) (int, int, error) {
	MID, errM := getMID(c)
	DPos, errDP := getID(c, "DPos", database.ERR_DAY_NOT_FOUND)

	if errM != nil {
		return MID, DPos, errDP
	} else if errDP != nil {
		return MID, DPos, errDP
	} else {
		return MID, DPos, nil
	}
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
	var MID int

	name := c.R.FormValue("name")
	meals, _ := strconv.Atoi(c.R.FormValue("meals"))
	days := strings.Split(c.R.FormValue("days"), "\n")

	if MID, err = c.U.Menus().New(name, days, meals); err == nil {
		utils.Redirect(c, "/menus/"+strconv.Itoa(MID)+"/edit")
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

	if MID, err = getMID(c); err == nil {
		if err = c.U.Menus().SetName(MID, c.R.FormValue("name")); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID))
		}
	}

	return
}

func GetMenuEditDay(c *utils.Context) (err error) {
	var MID, DPos int
	var day database.Day

	if MID, DPos, err = getDPos(c); err == nil {
		if day, err = c.U.Menus().GetDay(MID, DPos); err == nil {
			utils.RenderComponent(c, components.MenuEditDay(day))
		}
	}

	return
}

func PostMenuEditDayMeals(c *utils.Context) (err error) {
	var MID, DPos int
	var keys, meals []string
	c.R.ParseForm()

	if MID, DPos, err = getDPos(c); err == nil {
		for key, _ := range c.R.PostForm {
			if strings.HasPrefix(key, "meal-") {
				keys = append(keys, key)
			}
		}

		slices.Sort(keys)
		for _, key := range keys {
			meals = append(meals, c.R.FormValue(key))

		}

		if err = c.U.Menus().SetDayMeals(MID, DPos, meals); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID)+"/edit")
		}
	}

	return
}

func PostMenuEditDayName(c *utils.Context) (err error) {
	var MID, DPos int

	name := c.R.FormValue("name")

	if MID, DPos, err = getDPos(c); err == nil {
		if err = c.U.Menus().SetDayName(MID, DPos, name); err == nil {
			utils.Redirect(c, "/menus/"+strconv.Itoa(MID)+"/edit")
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

	if MID, err = getMID(c); err == nil {
		if _, err = c.U.Menus().Duplicate(MID); err == nil {
			utils.Redirect(c, "/menus")
		}
	}

	return
}
