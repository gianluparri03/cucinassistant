package database

import (
	"reflect"
	"strconv"
	"testing"
)

var testingEntriesN int = 0

func (sl ShoppingList) generate() (entry Entry) {
	testingEntriesN++

	entry.EID = testingEntriesN
	entry.Name = "entry-" + strconv.Itoa(testingEntriesN)
	sl.Append(entry.Name)

	if testingEntriesN%2 > 0 {
		sl.Toggle(testingEntriesN)
		entry.Marked = true
	}

	return
}

func TestShoppingListAppend(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	entry1 := Entry{EID: testingEntriesN + 1, Name: "appended-1"}
	entry2 := Entry{EID: testingEntriesN + 2, Name: "appended-2"}
	entry3 := Entry{EID: testingEntriesN + 3, Name: "appended-3"}
	testingEntriesN += 4

	names := []string{entry1.Name, entry2.Name, entry3.Name, entry2.Name}
	list := []Entry{entry1, entry2, entry3}

	type data struct {
		S     ShoppingList
		Names []string

		ExpectedErr  error
		ExpectedList []Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.S.Append(d.Names...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				list, _ := d.S.GetAll()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user appended entries",
				data{S: unknownUser.ShoppingList(), Names: names, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{S: s, Names: names, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestShoppingListClear(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	e1 := s.generate()
	e2 := s.generate()
	e3 := s.generate()

	var list []Entry
	for _, e := range []Entry{e1, e2, e3} {
		if !e.Marked {
			list = append(list, e)
		}
	}

	otherU, _ := getTestingUser(t)
	otherS := otherU.ShoppingList()
	otherEntry := otherS.generate()

	type data struct {
		S ShoppingList

		ExpectedErr  error
		ExpectedList []Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.S.Clear(); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				list, _ := d.S.GetAll()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user cleared shopping list",
				data{S: unknownUser.ShoppingList(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{S: s, ExpectedList: list},
			},
		},
	}.Run(t)

	list, _ = otherS.GetAll()
	if !reflect.DeepEqual(list, []Entry{otherEntry}) {
		t.Errorf("cleared shopping list of everyone %v", list)
	}
}

func TestShoppingListEdit(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	entry1 := s.generate()
	entry2 := s.generate()

	otherU, _ := getTestingUser(t)
	otherS := otherU.ShoppingList()

	type data struct {
		S       ShoppingList
		EID     int
		NewName string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.S.Edit(d.EID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.S.GetOne(d.EID)
				if entry.Name != d.NewName {
					t.Errorf("%s: name not changed", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited entry",
				data{S: otherS, EID: entry1.EID, NewName: entry1.Name + "+", ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"edited unknown entry",
				data{S: s, NewName: entry1.Name + "+", ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"duplicated entry",
				data{S: s, EID: entry1.EID, NewName: entry2.Name, ExpectedErr: ERR_ENTRY_DUPLICATED},
			},
			{
				"(same)",
				data{S: s, EID: entry1.EID, NewName: entry1.Name},
			},
			{
				"(different)",
				data{S: s, EID: entry1.EID, NewName: entry1.Name + "+"},
			},
		},
	}.Run(t)
}

func TestShoppingListGetAll(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	entry1 := s.generate()
	entry2 := s.generate()
	entry3 := s.generate()
	list := []Entry{entry1, entry2, entry3}

	otherU, _ := getTestingUser(t)
	otherS := otherU.ShoppingList()

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.ShoppingList().generate()

	type data struct {
		S ShoppingList

		ExpectedErr  error
		ExpectedList []Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			list, err := d.S.GetAll()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(list, d.ExpectedList) {
				t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
			}
		},

		Cases: []testCase[data]{
			{
				"got entries of unknown user",
				data{S: unknownUser.ShoppingList(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{S: otherS},
			},
			{
				"(filled)",
				data{S: s, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestShoppingListGetOne(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	entry1 := s.generate()
	entry2 := s.generate()

	otherU, _ := getTestingUser(t)
	otherS := otherU.ShoppingList()

	type data struct {
		S   ShoppingList
		EID int

		ExpectedErr   error
		ExpectedEntry Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			entry, err := d.S.GetOne(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(entry, d.ExpectedEntry) {
				t.Errorf("%s: expected entry <%v>, got <%v>", msg, d.ExpectedEntry, entry)
			}
		},

		Cases: []testCase[data]{
			{
				"other user retrieved entry",
				data{S: otherS, EID: entry1.EID, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"got data of unknown entry",
				data{S: s, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"(marked)",
				data{S: s, EID: entry1.EID, ExpectedEntry: entry1},
			},
			{
				"(unmarked)",
				data{S: s, EID: entry2.EID, ExpectedEntry: entry2},
			},
		},
	}.Run(t)
}

func TestShoppingListToggle(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.ShoppingList()

	entry := s.generate()
	before := entry.Marked

	otherU, _ := getTestingUser(t)
	otherS := otherU.ShoppingList()

	type data struct {
		S   ShoppingList
		EID int

		ExpectedErr    error
		ExpectedStatus bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.S.Toggle(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.S.GetOne(d.EID)
				if entry.Marked != d.ExpectedStatus {
					t.Errorf("%s: expected status <%v>, got <%v>", msg, d.ExpectedStatus, entry.Marked)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user toggled entry",
				data{S: otherS, EID: entry.EID, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"toggled unknown entry",
				data{S: s, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"(first time)",
				data{S: s, EID: entry.EID, ExpectedStatus: !before},
			},
			{
				"(second time)",
				data{S: s, EID: entry.EID, ExpectedStatus: before},
			},
		},
	}.Run(t)
}
