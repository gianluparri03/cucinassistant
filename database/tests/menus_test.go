package tests

import (
	"reflect"
	"testing"

	"cucinassistant/database"
)

func TestGetMenus(t *testing.T) {
	user, _ := GetTestingUser(t)
	m1, _ := user.NewMenu()
	m2, _ := user.NewMenu()
	m2, _ = user.ReplaceMenu(m2.MID, "newName", m2.Meals)

	otherUser, _ := GetTestingUser(t)

	otherOtherUser, _ := GetTestingUser(t)
	otherOtherUser.NewMenu()

	type data struct {
		User *database.User

		ExpectedErr   error
		ExpectedMenus []*database.Menu
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menus, err := d.User.GetMenus()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(menus, d.ExpectedMenus) {
				t.Errorf("%s: expected menus <%v>, got <%v>", msg, d.ExpectedMenus, menus)
			}
		},

		Cases: []TestCase[data]{
			{
				"got menus of unknown user",
				data{User: unknownUser, ExpectedErr: database.ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{User: otherUser, ExpectedMenus: []*database.Menu{}},
			},
			{
				"(filled)",
				data{User: user, ExpectedMenus: []*database.Menu{m1, m2}},
			},
		},
	}.Run(t)
}

func TestGetMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()
	menu, _ = user.ReplaceMenu(menu.MID, "newName", [14]string{"a", "b", "c"})

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		MID  int

		ExpectedErr  error
		ExpectedMenu *database.Menu
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menu, err := d.User.GetMenu(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(menu, d.ExpectedMenu) {
				t.Errorf("%s: expected menu <%v>, got <%v>", msg, d.ExpectedMenu, menu)
			}
		},

		Cases: []TestCase[data]{
			{
				"got data of unknown menu",
				data{User: user, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"other user retrieved menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, ExpectedMenu: menu},
			},
		},
	}.Run(t)
}

func TestNewMenu(t *testing.T) {
	user, _ := GetTestingUser(t)

	type data struct {
		User *database.User

		ExpectedErr error
		ExpectedMN  int
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if _, err := d.User.NewMenu(); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if menus, _ := d.User.GetMenus(); len(menus) != d.ExpectedMN {
				t.Errorf("%v, wrong number of menus", msg)
			}
		},

		Cases: []TestCase[data]{
			{
				"unknown user created menu",
				data{User: unknownUser, ExpectedErr: database.ERR_USER_UNKNOWN},
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

func TestReplaceMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()
	newName := "newName"
	newMeals := [14]string{"a", "b", "c"}

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User     *database.User
		MID      int
		NewName  string
		NewMeals [14]string

		ExpectedErr error
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.ReplaceMenu(d.MID, d.NewName, d.NewMeals)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				menu, _ := d.User.GetMenu(d.MID)
				expected := &database.Menu{MID: d.MID, Name: d.NewName, Meals: d.NewMeals}
				if !reflect.DeepEqual(menu, expected) {
					t.Errorf("%v, changes not saved", msg)
				} else if !reflect.DeepEqual(menu, got) {
					t.Errorf("%v, wrong returned values", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user replaced menu",
				data{User: otherUser, MID: menu.MID, NewName: newName, NewMeals: newMeals, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, NewName: newName, NewMeals: newMeals},
			},
		},
	}.Run(t)
}

func TestDeleteMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		MID  int

		ExpectedErr error
		ShouldExist bool
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.DeleteMenu(d.MID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			menu, _ = user.GetMenu(d.MID)
			if !d.ShouldExist && menu != nil {
				t.Errorf("%s, menu wasn't deleted", msg)
			} else if d.ShouldExist && menu == nil {
				t.Errorf("%s, menu was deleted anyway", msg)
			}
		},

		Cases: []TestCase[data]{
			{
				"other user deleted menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: database.ERR_MENU_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown menu",
				data{User: user, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID},
			},
		},
	}.Run(t)
}

func TestDuplicateMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()
	menu, _ = user.ReplaceMenu(menu.MID, "newName", [14]string{"a", "b", "c"})

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		MID  int

		ExpectedErr   error
		ExpectedMenus []*database.Menu
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.DuplicateMenu(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				srcMenu, _ := d.User.GetMenu(d.MID)
				dstMenu, _ := d.User.GetMenu(got.MID)
				if srcMenu.Meals != dstMenu.Meals {
					t.Errorf("%v, changes not saved", msg)
				} else if !reflect.DeepEqual(dstMenu, got) {
					t.Errorf("%v, wrong returned values", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user duplicated menu",
				data{User: otherUser, MID: menu.MID, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"duplicated unknown menu",
				data{User: user, ExpectedErr: database.ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID},
			},
		},
	}.Run(t)
}
