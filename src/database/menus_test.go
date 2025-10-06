package database

import (
	"reflect"
	"testing"
)

func TestMenusDelete(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("menu", []string{"a", "b", "c"}, 1)
	menu, _ := m.GetOne(MID)

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M   Menus
		MID int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.M.Delete(d.MID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			menu, _ = m.GetOne(d.MID)
			if !d.ShouldExist && menu.MID != 0 {
				t.Errorf("%s, menu wasn't deleted", msg)
			} else if d.ShouldExist && menu.MID == 0 {
				t.Errorf("%s, menu was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted menu",
				data{M: otherM, MID: MID, ExpectedErr: ERR_MENU_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID},
			},
		},
	}.Run(t)
}

func TestMenusDuplicate(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("m", []string{"d0", "d1", "d2"}, 3)
	menu, _ := m.GetOne(MID)

	meals0 := []string{"", "meal-0-1"}
	meals1 := []string{"", "", "", "meal-1-3"}
	meals2 := []string{}

	m.SetDayMeals(MID, 1, meals1)
	m.SetDayMeals(MID, 2, meals2)
	m.SetDayMeals(MID, 0, meals0)

	menu.Days[1].Meals = meals1
	menu.Days[2].Meals = meals2
	menu.Days[0].Meals = meals0

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M   Menus
		MID int

		ExpectedErr  error
		ExpectedMenu Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			mid, err := d.M.Duplicate(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				dst, _ := d.M.GetOne(mid)

				d.ExpectedMenu.MID = dst.MID
				for i, _ := range d.ExpectedMenu.Days {
					d.ExpectedMenu.Days[i].MID = dst.MID
				}

				if !reflect.DeepEqual(dst, d.ExpectedMenu) {
					t.Errorf("%s: expected menu <%v>, got <%v>", msg, d.ExpectedMenu, dst)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"duplicated unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"other user duplicated menu",
				data{M: otherM, MID: MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID, ExpectedMenu: menu},
			},
		},
	}.Run(t)
}

func TestMenusGetAll(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID1, _ := m.New("m1", []string{"d10", "d11", "d12"}, 3)
	MID2, _ := m.New("m2", []string{"d20", "d21", "d22"}, 2)

	menu1, _ := m.GetOne(MID1)
	menu2, _ := m.GetOne(MID2)

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.Menus().New("m0", []string{"d02"}, 5)

	type data struct {
		M Menus

		ExpectedErr   error
		ExpectedMenus []Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menus, err := d.M.GetAll()

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
				"(empty)",
				data{M: otherM},
			},
			{
				"(filled)",
				data{M: m, ExpectedMenus: []Menu{menu1, menu2}},
			},
		},
	}.Run(t)
}

func TestMenusGetDay(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("m", []string{"d0", "d1", "d2"}, 3)
	menu, _ := m.GetOne(MID)

	DPos := 2
	meals := []string{"meal-2-0", ""}
	m.SetDayMeals(MID, DPos, meals)
	m.SetDayName(MID, DPos, "d2")
	menu.Days[DPos].Meals = meals

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M    Menus
		MID  int
		DPos int

		ExpectedErr error
		ExpectedDay Day
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			day, err := d.M.GetDay(d.MID, d.DPos)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(day, d.ExpectedDay) {
				t.Errorf("%s: expected day <%v>, got <%v>", msg, d.ExpectedDay, day)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown day",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"other user retrieved day",
				data{M: otherM, MID: MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID, DPos: DPos, ExpectedDay: menu.Days[DPos]},
			},
		},
	}.Run(t)
}

func TestMenusGetOne(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("m", []string{"d0", "d1", "d2"}, 3)
	menu, _ := m.GetOne(MID)

	meals0 := []string{"", "meal-0-1"}
	meals1 := []string{"", "", "", "meal-1-3"}
	meals2 := []string{}

	m.SetDayMeals(MID, 1, meals1)
	m.SetDayMeals(MID, 2, meals2)
	m.SetDayMeals(MID, 0, meals0)

	menu.Days[1].Meals = meals1
	menu.Days[2].Meals = meals2
	menu.Days[0].Meals = meals0

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M   Menus
		MID int

		ExpectedErr  error
		ExpectedMenu Menu
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			menu, err := d.M.GetOne(d.MID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(menu, d.ExpectedMenu) {
				t.Errorf("%s: expected menu <%v>, got <%v>", msg, d.ExpectedMenu, menu)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"other user retrieved menu",
				data{M: otherM, MID: MID, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID, ExpectedMenu: menu},
			},
		},
	}.Run(t)
}

func TestMenusNew(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	dnames := []string{"d0", "d1", "d2", "d3"}

	type data struct {
		M         Menus
		DaysNames []string
		MealsN    int

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			name := "testMenu"

			if MID, err := d.M.New(name, d.DaysNames, d.MealsN); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				menu := Menu{MID: MID, Name: name}
				meals := make([]string, d.MealsN)
				for dpos, dname := range d.DaysNames {
					menu.Days = append(menu.Days, Day{MID: MID, Name: dname, Position: dpos, Meals: meals})
				}

				got, _ := d.M.GetOne(MID)

				if !reflect.DeepEqual(got, menu) {
					t.Errorf("%s: expected menu <%v>, got <%v>", msg, menu, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created menu",
				data{M: unknownUser.Menus(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"mealsN negative",
				data{M: m, DaysNames: dnames, MealsN: -1, ExpectedErr: ERR_MEALS_NEGATIVE},
			},
			{
				"(days>0, mealsN>0)",
				data{M: m, DaysNames: dnames, MealsN: 2},
			},
			{
				"(days>0, mealsN=0)",
				data{M: m, DaysNames: dnames, MealsN: 0},
			},
			{
				"(days=0, mealsN>0)",
				data{M: m, MealsN: 2},
			},
		},
	}.Run(t)
}

func TestMenuSetDayMeals(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("test", []string{"d0", "d1", "d2"}, 3)

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M     Menus
		MID   int
		Day   int
		Meals []string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			before, _ := d.M.GetOne(d.MID)

			if err := d.M.SetDayMeals(d.MID, d.Day, d.Meals); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Days[d.Day].Meals = d.Meals
				got, _ := d.M.GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: meals not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user set meal",
				data{M: otherM, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set meal of unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set meal of unknown day",
				data{M: m, MID: MID, Day: -1, ExpectedErr: ERR_DAY_NOT_FOUND},
			},
			{
				"(mealsN=)",
				data{M: m, MID: MID, Day: 1, Meals: []string{"meal0", "meal1", "meal2"}},
			},
			{
				"(mealsN<)",
				data{M: m, MID: MID, Day: 2, Meals: []string{}},
			},
			{
				"(mealsN>)",
				data{M: m, MID: MID, Day: 0, Meals: []string{"meal0", "meal1", "meal2", "meal3"}},
			},
		},
	}.Run(t)
}

func TestMenuSetDayName(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("test", []string{"d0", "d1", "d2"}, 3)

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M    Menus
		MID  int
		Day  int
		Name string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			before, _ := d.M.GetOne(d.MID)

			if err := d.M.SetDayName(d.MID, d.Day, d.Name); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Days[d.Day].Name = d.Name
				got, _ := d.M.GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: name not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user set day name",
				data{M: otherM, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set day name of unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set name of unknown day",
				data{M: m, MID: MID, Day: -1, ExpectedErr: ERR_DAY_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID, Day: 0, Name: "c"},
			},
		},
	}.Run(t)
}

func TestMenuSetName(t *testing.T) {
	u, _ := getTestingUser(t)
	m := u.Menus()

	MID, _ := m.New("test", []string{"d0", "d1", "d2"}, 3)

	otherU, _ := getTestingUser(t)
	otherM := otherU.Menus()

	type data struct {
		M   Menus
		MID int

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			newName := "newName"
			before, _ := d.M.GetOne(d.MID)

			if err := d.M.SetName(d.MID, newName); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				before.Name = newName
				got, _ := d.M.GetOne(d.MID)
				if !reflect.DeepEqual(before, got) {
					t.Errorf("%s: name not saved: expected <%v> got <%v>", msg, before, got)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user set menu name",
				data{M: otherM, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"set name of unknown menu",
				data{M: m, ExpectedErr: ERR_MENU_NOT_FOUND},
			},
			{
				"",
				data{M: m, MID: MID},
			},
		},
	}.Run(t)
}
