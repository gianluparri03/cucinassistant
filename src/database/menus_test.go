package database

import (
	"reflect"
	"testing"
)

func TestMenusGetAll(t *testing.T) {
	user, _ := getTestingUser(t)
	m1, _ := user.Menus().New("m1")
	m2, _ := user.Menus().New("m2")

	otherUser, _ := getTestingUser(t)

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.Menus().New("m0")

	type data struct {
		User User

		ExpectedErr   error
		ExpectedMenus []Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menus, err := d.User.Menus().GetAll()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(menus, d.ExpectedMenus) {
				t.Errorf("%s: expected menus <%v>, got <%v>", msg, d.ExpectedMenus, menus)
			}
		},

		Cases: []testCase[data]{
			{
				"got menus of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{User: otherUser},
			},
			{
				"(filled)",
				data{User: user, ExpectedMenus: []Menu{m1, m2}},
			},
		},
	}.Run(t)
}

func TestMenusGetOne(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("m")
	menu, _ = user.Menus().Replace(menu.MID, "m", [14]string{"a", "b", "c"})

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		MID  int

		ExpectedErr  error
		ExpectedMenu Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menu, err := d.User.Menus().GetOne(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(menu, d.ExpectedMenu) {
				t.Errorf("%s: expected menu <%v>, got <%v>", msg, d.ExpectedMenu, menu)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"other user retrieved menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, ExpectedMenu: menu},
			},
		},
	}.Run(t)
}

func TestMenusNew(t *testing.T) {
	user, _ := getTestingUser(t)

	type data struct {
		User User
		Name string

		ExpectedErr error
		ExpectedMN  int
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if m, err := d.User.Menus().New("name"); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil && m.Name != "name" {
				t.Errorf("%s: expected name <test>, got <%v>", msg, m.Name)
			}

			if menus, _ := d.User.Menus().GetAll(); len(menus) != d.ExpectedMN {
				t.Errorf("%v, wrong number of menus", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created menu",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, ExpectedMN: 1},
			},
			{
				"",
				data{User: user, ExpectedMN: 2},
			},
		},
	}.Run(t)
}

func TestMenusReplace(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("oldName")
	newName := "newName"
	newMeals := [14]string{"a", "b", "c"}

	otherUser, _ := getTestingUser(t)

	type data struct {
		User     User
		MID      int
		NewName  string
		NewMeals [14]string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.Menus().Replace(d.MID, d.NewName, d.NewMeals)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				menu, _ := d.User.Menus().GetOne(d.MID)
				expected := Menu{MID: d.MID, Name: d.NewName, Meals: d.NewMeals}
				if !reflect.DeepEqual(menu, expected) {
					t.Errorf("%v, changes not saved", msg)
				} else if !reflect.DeepEqual(menu, got) {
					t.Errorf("%v, new menu badly returned", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user replaced menu",
				data{User: otherUser, MID: menu.MID, NewName: newName, NewMeals: newMeals, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"replaced unknown menu",
				data{User: user, NewName: newName, NewMeals: newMeals, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, NewName: newName, NewMeals: newMeals},
			},
		},
	}.Run(t)
}

func TestMenusDelete(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("")

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		MID  int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.Menus().Delete(d.MID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			menu, _ = user.Menus().GetOne(d.MID)
			if !d.ShouldExist && menu.MID != 0 {
				t.Errorf("%s, menu wasn't deleted", msg)
			} else if d.ShouldExist && menu.MID == 0 {
				t.Errorf("%s, menu was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: ERR_MENU_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID},
			},
		},
	}.Run(t)
}

func TestMenusDuplicate(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("")
	menu, _ = user.Menus().Replace(menu.MID, "name", [14]string{"a", "b", "c"})

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		MID  int

		ExpectedErr   error
		ExpectedMenus []*Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.Menus().Duplicate(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				srcMenu, _ := d.User.Menus().GetOne(d.MID)
				dstMenu, _ := d.User.Menus().GetOne(got.MID)
				if srcMenu.Meals != dstMenu.Meals {
					t.Errorf("%v, changes not saved", msg)
				} else if !reflect.DeepEqual(dstMenu, got) {
					t.Errorf("%v, wrong returned values", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user duplicated menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"duplicated unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID},
			},
		},
	}.Run(t)
}
