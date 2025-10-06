package database

import (
	"reflect"
	"testing"
)

func TestMenusDelete(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("menu", []string{"a", "b", "c"}, 1)

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

func TestMenusGetAll(t *testing.T) {
	user, _ := getTestingUser(t)
	m1, _ := user.Menus().New("m1", []string{"d10", "d11", "d12"}, 3)
	m2, _ := user.Menus().New("m2", []string{"d20", "d21", "d22"}, 2)

	otherUser, _ := getTestingUser(t)

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.Menus().New("m0", []string{"d02"}, 5)

	type data struct {
		User User

		ExpectedErr   error
		ExpectedMenus []Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menus, err := d.User.Menus().GetAll()

			expected := d.ExpectedMenus
			for i, _ := range expected {
				expected[i].Days = nil
			}

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

func TestMenusGetDay(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("m", []string{"d0", "d1", "d2"}, 3)

	meals := []string{"meal-2-0", ""}
	user.Menus().SetDayMeals(menu.MID, 2, meals)
	user.Menus().SetDayName(menu.MID, 2, "d2")
	day := menu.Days[2]
	day.Meals = meals

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		MID  int
		DPos int

		ExpectedErr error
		ExpectedDay Day
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			day, err := d.User.Menus().GetDay(d.MID, d.DPos)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(day, d.ExpectedDay) {
				t.Errorf("%s: expected day <%v>, got <%v>", msg, d.ExpectedDay, day)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown day",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"other user retrieved day",
				data{User: otherUser, MID: menu.MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, DPos: 2, ExpectedDay: day},
			},
		},
	}.Run(t)
}

func TestMenusGetOne(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("m", []string{"d0", "d1", "d2"}, 3)

	meals0 := []string{"", "meal-0-1"}
	meals1 := []string{"", "", "", "meal-1-3"}
	meals2 := []string{"meal-2-0", ""}

	user.Menus().SetDayMeals(menu.MID, 1, meals1)
	user.Menus().SetDayMeals(menu.MID, 2, meals2)
	user.Menus().SetDayMeals(menu.MID, 0, meals0)

	menu.Days[1].Meals = meals1
	menu.Days[2].Meals = meals2
	menu.Days[0].Meals = meals0

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

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			name := "testMenu"
			daysNames := []string{"day1", "day2", "day3"}
			mealsN := 1

			if m, err := d.User.Menus().New(name, daysNames, mealsN); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				menu := Menu{
					MID:  m.MID,
					Name: name,
					Days: []Day{
						Day{MID: m.MID, Name: "day1", Position: 0, Meals: make([]string, mealsN)},
						Day{MID: m.MID, Name: "day2", Position: 1, Meals: make([]string, mealsN)},
						Day{MID: m.MID, Name: "day3", Position: 2, Meals: make([]string, mealsN)},
					},
				}

				if !reflect.DeepEqual(m, menu) {
					t.Errorf("%s: expected menu <%v>, got <%v>", msg, menu, m)
				} else if m2, err := d.User.Menus().GetOne(m.MID); err != nil {
					t.Errorf("%s: menu not saved (error)", msg)
				} else if !reflect.DeepEqual(m2, menu) {
					t.Errorf("%s: menu not saved (differ): expected <%v> got <%v>", msg, menu, m2)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created menu",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user},
			},
		},
	}.Run(t)
}

func TestMenuSetDayMeals(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("test", []string{"d0", "d1", "d2"}, 3)

	type data struct {
		User  User
		MID   int
		Day   int
		Meals []string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			before, _ := d.User.Menus().GetOne(d.MID)

			if err := d.User.Menus().SetDayMeals(d.MID, d.Day, d.Meals); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Days[d.Day].Meals = d.Meals
				got, _ := d.User.Menus().GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: meals not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user set meal",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"set meal of unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set meal of unknown day",
				data{User: user, MID: menu.MID, Day: -1, ExpectedErr: ERR_DAY_NOT_FOUND},
			},
			{
				"(mealsN=)",
				data{User: user, MID: menu.MID, Day: 1, Meals: []string{"meal0", "meal1", "meal2"}},
			},
			{
				"(mealsN<)",
				data{User: user, MID: menu.MID, Day: 2, Meals: []string{}},
			},
			{
				"(mealsN>)",
				data{User: user, MID: menu.MID, Day: 0, Meals: []string{"meal0", "meal1", "meal2", "meal3"}},
			},
		},
	}.Run(t)
}

func TestMenuSetDayName(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("test", []string{"d0", "d1", "d2"}, 3)

	type data struct {
		User User
		MID  int
		Day  int
		Name string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			before, _ := d.User.Menus().GetOne(d.MID)

			if err := d.User.Menus().SetDayName(d.MID, d.Day, d.Name); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Days[d.Day].Name = d.Name
				got, _ := d.User.Menus().GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: name not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user set day name",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"set day name of unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set name of unknown day",
				data{User: user, MID: menu.MID, Day: -1, ExpectedErr: ERR_DAY_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID, Day: 0, Name: "c"},
			},
		},
	}.Run(t)
}

func TestMenuSetName(t *testing.T) {
	user, _ := getTestingUser(t)
	menu, _ := user.Menus().New("test", []string{"d0", "d1", "d2"}, 3)

	type data struct {
		User User
		MID  int

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			newName := "newName"
			before, _ := d.User.Menus().GetOne(d.MID)

			if err := d.User.Menus().SetName(d.MID, newName); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Name = newName
				got, _ := d.User.Menus().GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: name not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user set menu name",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"set name of unknown menu",
				data{User: user, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{User: user, MID: menu.MID},
			},
		},
	}.Run(t)
}
