package database

import (
	"strconv"
)

var testingEntriesN int = 0

// generateEntry adds an entry to an user's shopping list
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
