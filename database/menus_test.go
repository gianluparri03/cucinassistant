package database

import (
	"reflect"
	"testing"
)

func TestGetMenus(t *testing.T) {
	userEmpty, _ := GetTestingUser(t)
	userFilled, _ := GetTestingUser(t)

	m1, _ := userFilled.NewMenu()
	m2, _ := userFilled.NewMenu()

	m2.Name = "newName"
	userFilled.ReplaceMenu(m2.MID, m2)

	TestSuite[Pair[[]*Menu, error]]{
		Target: func(tc *TestCase[Pair[[]*Menu, error]]) Pair[[]*Menu, error] {
			menus, err := tc.User.GetMenus()
			return Pair[[]*Menu, error]{menus, err}
		},

		Cases: []TestCase[Pair[[]*Menu, error]]{
			{
				Description: "got menus of unknown user",
				User:        &User{UID: 0},
				Expected:    Pair[[]*Menu, error]{nil, ERR_USER_UNKNOWN},
			},
			{
				Description: "wrong menus (empty)",
				User:        userEmpty,
				Expected:    Pair[[]*Menu, error]{[]*Menu{}, nil},
			},
			{
				Description: "wrong menus (filled)",
				User:        userFilled,
				Expected:    Pair[[]*Menu, error]{[]*Menu{m1, m2}, nil},
			},
		},
	}.Run(t)
}

func TestGetMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()
	menu.Name = "newName"
	menu.Meals = [14]string{"a", "b", "c"}
	user.ReplaceMenu(menu.MID, menu)

	TestSuite[Pair[*Menu, error]]{
		Target: func(tc *TestCase[Pair[*Menu, error]]) Pair[*Menu, error] {
			menu, err := tc.User.GetMenu(tc.Data["MID"].(int))
			return Pair[*Menu, error]{menu, err}
		},

		Cases: []TestCase[Pair[*Menu, error]]{
			{
				Description: "unknown user retrieved menu",
				User:        &User{UID: 0},
				Expected:    Pair[*Menu, error]{nil, ERR_USER_UNKNOWN},
				Data:        map[string]any{"MID": 0},
			},
			{
				Description: "got data of unknown menu",
				User:        user,
				Expected:    Pair[*Menu, error]{nil, ERR_MENU_NOT_FOUND},
				Data:        map[string]any{"MID": 0},
			},
			{
				Description: "wrong menu data",
				User:        user,
				Expected:    Pair[*Menu, error]{menu, nil},
				Data:        map[string]any{"MID": menu.MID},
			},
		},
	}.Run(t)
}

func TestNewMenu(t *testing.T) {
	user, _ := GetTestingUser(t)

	TestSuite[Pair[*Menu, error]]{
		Target: func(tc *TestCase[Pair[*Menu, error]]) Pair[*Menu, error] {
			menu, err := tc.User.NewMenu()
			return Pair[*Menu, error]{menu, err}
		},

		PostCheck: func(t *testing.T, tc *TestCase[Pair[*Menu, error]]) {
			expected := tc.Data["MenusN"].(int)
			menus, _ := tc.User.GetMenus()
			if len(menus) != expected {
				t.Errorf("%v, number of menus does not match: expected <%v> got <%v>", tc.Description, expected, len(menus))
			}
		},

		Cases: []TestCase[Pair[*Menu, error]]{
			{
				Description: "unknown user created menu",
				User:        &User{UID: 0},
				Expected:    Pair[*Menu, error]{nil, ERR_USER_UNKNOWN},
				Data:        map[string]any{"MenusN": 0},
			},
			{
				Description: "could not create menu",
				User:        user,
				Expected:    Pair[*Menu, error]{&Menu{MID: 1, Name: menuDefaultName}, nil},
				Data:        map[string]any{"MenusN": 1},
			},
		},
	}.Run(t)
}

func TestReplaceMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()

	newName := "newName"
	newMeals := [14]string{"a", "b", "c"}
	expected := &Menu{MID: menu.MID, Name: newName, Meals: newMeals}

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			newData := &Menu{Name: tc.Data["NewName"].(string), Meals: tc.Data["NewMeals"].([14]string)}
			return tc.User.ReplaceMenu(tc.Data["MID"].(int), newData)
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				expected := tc.Data["Menu"].(*Menu)
				got, _ := tc.User.GetMenu(tc.Data["MID"].(int))
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("%v, changes not saved: expected <%v> got <%v>", tc.Description, expected, got)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "unknown user replaced menu",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"MID": 0, "NewName": newName, "NewMeals": newMeals, "Menu": nil},
			},
			{
				Description: "could not replace menu",
				User:        user,
				Expected:    nil,
				Data:        map[string]any{"MID": menu.MID, "NewName": newName, "NewMeals": newMeals, "Menu": expected},
			},
		},
	}.Run(t)
}

func TestDeleteMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.DeleteMenu(tc.Data["MID"].(int))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				_, err := tc.User.GetMenu(tc.Data["MID"].(int))
				if err == nil {
					t.Errorf("%v, menu not deleted", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "unknown user deleted menu",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"MID": 0},
			},
			{
				Description: "deleted unknown menu",
				User:        user,
				Expected:    ERR_MENU_NOT_FOUND,
				Data:        map[string]any{"MID": 0},
			},
			{
				Description: "could not delete menu",
				User:        user,
				Expected:    nil,
				Data:        map[string]any{"MID": menu.MID},
			},
		},
	}.Run(t)
}

func TestDuplicateMenu(t *testing.T) {
	user, _ := GetTestingUser(t)
	menu, _ := user.NewMenu()
	menu.Name = "newName"
	menu.Meals = [14]string{"a", "b", "c"}
	user.ReplaceMenu(menu.MID, menu)

	TestSuite[Pair[*Menu, error]]{
		Target: func(tc *TestCase[Pair[*Menu, error]]) Pair[*Menu, error] {
			menu, err := tc.User.DuplicateMenu(tc.Data["MID"].(int))
			return Pair[*Menu, error]{menu, err}
		},

		PostCheck: func(t *testing.T, tc *TestCase[Pair[*Menu, error]]) {
			expected := tc.Data["MenusN"].(int)
			menus, _ := tc.User.GetMenus()
			if len(menus) != expected {
				t.Errorf("%v, number of menus does not match: expected <%v> got <%v>", tc.Description, expected, len(menus))
			}
		},

		Cases: []TestCase[Pair[*Menu, error]]{
			{
				Description: "unknown user duplicated menu",
				User:        &User{UID: 0},
				Expected:    Pair[*Menu, error]{nil, ERR_USER_UNKNOWN},
				Data:        map[string]any{"MID": 0, "MenusN": 0},
			},
			{
				Description: "duplicated unknown menu",
				User:        user,
				Expected:    Pair[*Menu, error]{nil, ERR_MENU_NOT_FOUND},
				Data:        map[string]any{"MID": 0, "MenusN": 1},
			},
			{
				Description: "could not duplicate menu",
				User:        user,
				Expected:    Pair[*Menu, error]{&Menu{MID: 2, Name: menu.Name + duplicatedMenuSuffix, Meals: menu.Meals}, nil},
				Data:        map[string]any{"MID": menu.MID, "MenusN": 2},
			},
		},
	}.Run(t)
}
