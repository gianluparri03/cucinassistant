package database

import (
	"reflect"
	"strconv"
	"testing"
)

var testingEntriesN int = 0

func (sl ShoppingList) generate() (entry Entry) {
	testingEntriesN++
	name := "entry-" + strconv.Itoa(testingEntriesN)
	sl.Append(name)

	// Find its id
	list, _ := sl.GetAll()
	var EID int
	for EID, entry = range list {
		if entry.Name == name {
			break
		}
	}

	// Marks only the odd ones
	if testingEntriesN%2 > 0 {
		sl.Toggle(EID)
		entry.Marked = true
	}

	return
}

func TestShoppingListGetAll(t *testing.T) {
	user, _ := getTestingUser(t)
	entry1 := user.ShoppingList().generate()
	entry2 := user.ShoppingList().generate()
	entry3 := user.ShoppingList().generate()
	list := map[int]Entry{entry1.EID: entry1, entry2.EID: entry2, entry3.EID: entry3}

	otherUser, _ := getTestingUser(t)

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.ShoppingList().generate()

	type data struct {
		User User

		ExpectedErr  error
		ExpectedList map[int]Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			list, err := d.User.ShoppingList().GetAll()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(list, d.ExpectedList) {
				t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
			}
		},

		Cases: []testCase[data]{
			{
				"got entries of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN, ExpectedList: map[int]Entry{}},
			},
			{
				"(empty)",
				data{User: otherUser, ExpectedList: map[int]Entry{}},
			},
			{
				"(filled)",
				data{User: user, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestShoppingListGetOne(t *testing.T) {
	user, _ := getTestingUser(t)
	entry1 := user.ShoppingList().generate()
	entry2 := user.ShoppingList().generate()

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		EID  int

		ExpectedErr   error
		ExpectedEntry Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			entry, err := d.User.ShoppingList().GetOne(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(entry, d.ExpectedEntry) {
				t.Errorf("%s: expected entry <%v>, got <%v>", msg, d.ExpectedEntry, entry)
			}
		},

		Cases: []testCase[data]{
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

func TestShoppingListAppend(t *testing.T) {
	user, _ := getTestingUser(t)
	entry1 := user.ShoppingList().generate()
	entry2 := Entry{EID: testingEntriesN + 1, Name: "appended-2"}
	entry3 := Entry{EID: testingEntriesN + 2, Name: "appended-3"}
	testingEntriesN += 2

	names := []string{entry2.Name, entry3.Name, entry1.Name, entry2.Name}
	list := map[int]Entry{entry1.EID: entry1, entry2.EID: entry2, entry3.EID: entry3}

	type data struct {
		User  User
		Names []string

		ExpectedErr  error
		ExpectedList map[int]Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ShoppingList().Append(d.Names...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				list, _ := d.User.ShoppingList().GetAll()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user appended entries",
				data{User: unknownUser, Names: names, ExpectedErr: ERR_USER_UNKNOWN, ExpectedList: map[int]Entry{}},
			},
			{
				"",
				data{User: user, Names: names, ExpectedList: list},
			},
		},
	}.Run(t)
}

func TestShoppingListToggle(t *testing.T) {
	user, _ := getTestingUser(t)
	entry := user.ShoppingList().generate()

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		EID  int

		ExpectedErr    error
		ExpectedStatus bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ShoppingList().Toggle(d.EID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.User.ShoppingList().GetOne(d.EID)
				if entry.Marked != d.ExpectedStatus {
					t.Errorf("%s: expected status <%v>, got <%v>", msg, d.ExpectedStatus, entry.Marked)
				}
			}
		},

		Cases: []testCase[data]{
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

func TestShoppingListClear(t *testing.T) {
	user, _ := getTestingUser(t)
	user.ShoppingList().generate()
	entry := user.ShoppingList().generate()
	user.ShoppingList().generate()

	otherUser, _ := getTestingUser(t)
	otherEntry := otherUser.ShoppingList().generate()

	type data struct {
		User User

		ExpectedErr  error
		ExpectedList map[int]Entry
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.ShoppingList().Clear(); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				list, _ := d.User.ShoppingList().GetAll()
				if !reflect.DeepEqual(list, d.ExpectedList) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedList, list)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user cleared shopping list",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, ExpectedList: map[int]Entry{entry.EID: entry}},
			},
		},
	}.Run(t)

	list, _ := otherUser.ShoppingList().GetAll()
	if !reflect.DeepEqual(list, map[int]Entry{otherEntry.EID: otherEntry}) {
		t.Errorf("cleared shopping list of everyone")
	}
}

func TestShoppingListEdit(t *testing.T) {
	user, _ := getTestingUser(t)
	entry1 := user.ShoppingList().generate()
	entry2 := user.ShoppingList().generate()

	otherUser, _ := getTestingUser(t)

	type data struct {
		User    User
		EID     int
		NewName string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ShoppingList().Edit(d.EID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				entry, _ := d.User.ShoppingList().GetOne(d.EID)
				if entry.Name != d.NewName {
					t.Errorf("%s: name not changed", msg)
				}
			}
		},

		Cases: []testCase[data]{
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
