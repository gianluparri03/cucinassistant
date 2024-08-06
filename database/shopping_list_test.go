package database

import (
	"reflect"
	"testing"
)

func TestGetShoppingList(t *testing.T) {
	list := ShoppingList{
		1: Entry{EID: 1, Name: "Entry1"},
		2: Entry{EID: 2, Name: "Entry2", Marked: true},
		3: Entry{EID: 3, Name: "Entry3"},
	}

	userEmpty, _ := GetTestingUser(t)
	userFilled, _ := GetTestingUser(t)
	userFilled.AppendEntries(list[1].Name, list[2].Name, list[3].Name)
	userFilled.ToggleEntry(list[2].EID)

	TestSuite[Pair[ShoppingList, error]]{
		Target: func(tc *TestCase[Pair[ShoppingList, error]]) Pair[ShoppingList, error] {
			list, e := tc.User.GetShoppingList()
			return Pair[ShoppingList, error]{list, e}
		},

		Cases: []TestCase[Pair[ShoppingList, error]]{
			{
				Description: "got entries of inexistent user",
				User:        &User{UID: 0},
				Expected:    Pair[ShoppingList, error]{ShoppingList{}, ERR_USER_UNKNOWN},
			},
			{
				Description: "wrong user shopping list (empty)",
				User:        &userEmpty,
				Expected:    Pair[ShoppingList, error]{ShoppingList{}, nil},
			},
			{
				Description: "wrong user shopping list (filled)",
				User:        &userFilled,
				Expected:    Pair[ShoppingList, error]{list, nil},
			},
		},
	}.Run(t)
}

func TestGetEntry(t *testing.T) {
	user, _ := GetTestingUser(t)

	entryUnmarked := Entry{EID: 1, Name: "unmarked"}
	entryMarked := Entry{EID: 2, Name: "marked", Marked: true}
	user.AppendEntries(entryUnmarked.Name, entryMarked.Name)
	user.ToggleEntry(entryMarked.EID)

	TestSuite[Pair[Entry, error]]{
		Target: func(tc *TestCase[Pair[Entry, error]]) Pair[Entry, error] {
			en, er := tc.User.GetEntry(tc.Data["EID"].(int))
			return Pair[Entry, error]{en, er}
		},

		Cases: []TestCase[Pair[Entry, error]]{
			{
				Description: "got data of inexistent entry",
				User:        &user,
				Expected:    Pair[Entry, error]{Entry{}, ERR_ENTRY_NOT_FOUND},
				Data:        map[string]any{"EID": 0},
			},
			{
				Description: "wrong entry's data (unmarked)",
				User:        &user,
				Expected:    Pair[Entry, error]{entryUnmarked, nil},
				Data:        map[string]any{"EID": entryUnmarked.EID},
			},
			{
				Description: "wrong entry's data (marked)",
				User:        &user,
				Expected:    Pair[Entry, error]{entryMarked, nil},
				Data:        map[string]any{"EID": entryMarked.EID},
			},
		},
	}.Run(t)
}

func TestAppendEntries(t *testing.T) {
	list := ShoppingList{
		1: Entry{EID: 1, Name: "Entry1"},
		2: Entry{EID: 2, Name: "Entry2"},
		3: Entry{EID: 3, Name: "Entry3"},
	}
	names := []string{list[1].Name, list[2].Name, list[2].Name, list[1].Name, list[3].Name}

	user, _ := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.AppendEntries(tc.Data["Names"].([]string)...)
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			got, _ := tc.User.GetShoppingList()
			expected := tc.Data["List"].(ShoppingList)
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("%s, list does not match: expected <%v> got <%v>", tc.Description, expected, got)
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "inexistent user appended entries",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"Names": names, "List": ShoppingList{}},
			},
			{
				Description: "could not append entries",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"Names": names, "List": list},
			},
		},
	}.Run(t)
}

func TestToggleEntry(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry := Entry{EID: 1, Name: "name"}
	user.AppendEntries(entry.Name)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ToggleEntry(tc.Data["EID"].(int))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				entry, _ := tc.User.GetEntry(tc.Data["EID"].(int))
				shouldBe := tc.Data["ShouldBe"].(bool)
				if entry.Marked != shouldBe {
					t.Errorf("%s, wrong marked value: expected <%v> got <%v>", tc.Description, shouldBe, entry.Marked)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "toggled inexistent entry",
				User:        &user,
				Expected:    ERR_ENTRY_NOT_FOUND,
				Data:        map[string]any{"EID": 0},
			},
			{
				Description: "could not toggle entry (false->true)",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"EID": entry.EID, "ShouldBe": true},
			},
			{
				Description: "could not toggle entry (true->false)",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"EID": entry.EID, "ShouldBe": false},
			},
		},
	}.Run(t)
}

func TestCleanEntries(t *testing.T) {
	list := ShoppingList{
		1: Entry{EID: 1, Name: "Entry1"},
		2: Entry{EID: 2, Name: "Entry2", Marked: true},
		3: Entry{EID: 3, Name: "Entry3"},
	}

	user, _ := GetTestingUser(t)
	user.AppendEntries(list[1].Name, list[2].Name, list[3].Name)
	user.ToggleEntry(list[2].EID)

	delete(list, 2)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.CleanEntries()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				got, _ := tc.User.GetShoppingList()
				expected := tc.Data["List"].(ShoppingList)
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("%s, shopping list does not match: expected <%v> got <%v>", tc.Description, expected, got)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "inexistent user cleaned its entries",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
			},
			{
				Description: "could not clean entries",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"List": list},
			},
		},
	}.Run(t)
}

func TestEditEntry(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry1 := Entry{EID: 1, Name: "Entry1"}
	entry2 := Entry{EID: 2, Name: "Entry2"}
	user.AppendEntries(entry1.Name, entry2.Name)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.EditEntry(tc.Data["EID"].(int), tc.Data["NewName"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				entry, _ := tc.User.GetEntry(tc.Data["EID"].(int))
				newName := tc.Data["NewName"].(string)
				if entry.Name != newName {
					t.Errorf("%s, new name does not match: expected <%v> got <%v>", tc.Description, newName, entry.Name)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "edited unknown entry",
				User:        &user,
				Expected:    ERR_ENTRY_NOT_FOUND,
				Data:        map[string]any{"EID": 0, "NewName": ""},
			},
			{
				Description: "duplicated entry",
				User:        &user,
				Expected:    ERR_ENTRY_DUPLICATED,
				Data:        map[string]any{"EID": entry2.EID, "NewName": entry1.Name},
			},
			{
				Description: "could not keep entry's name",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"EID": entry2.EID, "NewName": entry2.Name},
			},
			{
				Description: "could not edit entry",
				User:        &user,
				Expected:    nil,
				Data:        map[string]any{"EID": entry2.EID, "NewName": entry2.Name + "+"},
			},
		},
	}.Run(t)
}
