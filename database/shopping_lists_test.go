package database

import (
	"reflect"
	"strconv"
	"testing"
)

var testingEntriesN int = 0

// denerateEntry adds an entry to an user's shopping list
func generateEntry(u *User) (entry Entry) {
	testingEntriesN++

	name := "entry-" + strconv.Itoa(testingEntriesN)
	u.AppendEntries(name)

	// Find its id
	list, _ := u.GetShoppingList()
	var EID int
	for EID, entry = range list {
		if entry.Name == name {
			break
		}
	}

	// Marks only the odd ones
	if testingEntriesN%2 > 0 {
		u.ToggleEntry(EID)
		entry.Marked = true
	}

	return
}

func TestGetShoppingList(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry1 := generateEntry(user)
	entry2 := generateEntry(user)
	entry3 := generateEntry(user)
	list := ShoppingList{entry1.EID: entry1, entry2.EID: entry2, entry3.EID: entry3}

	otherUser, _ := GetTestingUser(t)

	otherOtherUser, _ := GetTestingUser(t)
	generateEntry(otherOtherUser)

	type data struct {
		User *User

		ExpectedErr  error
		ExpectedList ShoppingList
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			list, err := d.User.GetShoppingList()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(list, d.ExpectedList) {
				t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
			}
		},

		Cases: []TestCase[data]{
			{
				"got entries of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN, ExpectedList: ShoppingList{}},
			},
			{
				"(empty)",
				data{User: otherUser, ExpectedList: ShoppingList{}},
			},
			{
				"(filled)",
				data{User: user, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestGetEntry(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry1 := generateEntry(user)
	entry2 := generateEntry(user)

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *User
		EID  int

		ExpectedErr   error
		ExpectedEntry Entry
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			entry, err := d.User.GetEntry(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(entry, d.ExpectedEntry) {
				t.Errorf("%s: expected entry <%v>, got <%v>", msg, d.ExpectedEntry, entry)
			}
		},

		Cases: []TestCase[data]{
			{
				"other user retrieved entry",
				data{User: otherUser, EID: entry1.EID, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"got data of unknown entry",
				data{User: user, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"(marked)",
				data{User: user, EID: entry1.EID, ExpectedEntry: entry1},
			},
			{
				"(unmarked)",
				data{User: user, EID: entry2.EID, ExpectedEntry: entry2},
			},
		},
	}.Run(t)
}

func TestAppendEntries(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry1 := generateEntry(user)
	entry2 := Entry{EID: testingEntriesN + 1, Name: "appended-2"}
	entry3 := Entry{EID: testingEntriesN + 2, Name: "appended-3"}
	testingEntriesN += 2

	names := []string{entry2.Name, entry3.Name, entry1.Name, entry2.Name}
	list := ShoppingList{entry1.EID: entry1, entry2.EID: entry2, entry3.EID: entry3}

	type data struct {
		User  *User
		Names []string

		ExpectedErr  error
		ExpectedList ShoppingList
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.AppendEntries(d.Names...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				list, _ := d.User.GetShoppingList()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"unknown user appended entries",
				data{User: unknownUser, Names: names, ExpectedErr: ERR_USER_UNKNOWN, ExpectedList: ShoppingList{}},
			},
			{
				"",
				data{User: user, Names: names, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestToggleEntry(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry := generateEntry(user)

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *User
		EID  int

		ExpectedErr    error
		ExpectedStatus bool
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ToggleEntry(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.User.GetEntry(d.EID)
				if entry.Marked != d.ExpectedStatus {
					t.Errorf("%s: expected status <%v>, got <%v>", msg, d.ExpectedStatus, entry.Marked)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user toggled entry",
				data{User: otherUser, EID: entry.EID, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"toggled unknown entry",
				data{User: user, ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"(marking)",
				data{User: user, EID: entry.EID, ExpectedStatus: true},
			},
			{
				"(unmarking)",
				data{User: user, EID: entry.EID, ExpectedStatus: false},
			},
		},
	}.Run(t)
}

func TestClearShoppingList(t *testing.T) {
	user, _ := GetTestingUser(t)
	generateEntry(user)
	entry := generateEntry(user)
	generateEntry(user)

	otherUser, _ := GetTestingUser(t)
	otherEntry := generateEntry(otherUser)

	type data struct {
		User *User

		ExpectedErr  error
		ExpectedList ShoppingList
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.ClearShoppingList(); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				list, _ := d.User.GetShoppingList()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"unknown user cleared shopping list",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, ExpectedList: ShoppingList{entry.EID: entry}},
			},
		},
	}.Run(t)

	list, _ := otherUser.GetShoppingList()
	if !reflect.DeepEqual(list, ShoppingList{otherEntry.EID: otherEntry}) {
		t.Errorf("cleared shopping list of everyone")
	}
}

func TestEditEntry(t *testing.T) {
	user, _ := GetTestingUser(t)
	entry1 := generateEntry(user)
	entry2 := generateEntry(user)

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User    *User
		EID     int
		NewName string

		ExpectedErr error
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.EditEntry(d.EID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.User.GetEntry(d.EID)
				if entry.Name != d.NewName {
					t.Errorf("%s: name not changed", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user edited entry",
				data{User: otherUser, EID: entry1.EID, NewName: entry1.Name + "+", ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"edited unknown entry",
				data{User: user, NewName: entry1.Name + "+", ExpectedErr: ERR_ENTRY_NOT_FOUND},
			},
			{
				"duplicated entry",
				data{User: user, EID: entry1.EID, NewName: entry2.Name, ExpectedErr: ERR_ENTRY_DUPLICATED},
			},
			{
				"(same)",
				data{User: user, EID: entry1.EID, NewName: entry1.Name},
			},
			{
				"(different)",
				data{User: user, EID: entry1.EID, NewName: entry1.Name + "+"},
			},
		},
	}.Run(t)
}
