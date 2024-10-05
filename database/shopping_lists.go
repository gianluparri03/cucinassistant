package database

import (
	"database/sql"
	"log/slog"
)

// Entry is an element of the shopping list
type Entry struct {
	// EID is the Entry ID
	EID int

	// Name is the name of the entry
	Name string

	// Marked indicates if the checkbox has been checked
	Marked bool
}

// ShoppingList is an alias for a map of entries
type ShoppingList map[int]Entry

// GetShoppingList returns an user's shopping list
func (u *User) GetShoppingList() (sl ShoppingList, err error) {
	sl = ShoppingList{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = db.Query(`SELECT eid, name, marked FROM entries WHERE uid=$1;`, u.UID)
	if err != nil {
		slog.Error("while retrieving shopping list:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var e Entry
		rows.Scan(&e.EID, &e.Name, &e.Marked)
		sl[e.EID] = e
	}

	// If no entries have been found, makes sure the user exists
	if len(sl) == 0 {
		_, err = GetUser("UID", u.UID)
	}

	return
}

// GetEntry returns a single entry of the shopping list
func (u *User) GetEntry(EID int) (e Entry, err error) {
	// Fetches them
	err = db.QueryRow(`SELECT eid, name, marked FROM entries WHERE uid=$1 AND eid=$2;`, u.UID, EID).Scan(&e.EID, &e.Name, &e.Marked)
	if err != nil {
		err = handleNoRowsError(err, u.UID, ERR_ENTRY_NOT_FOUND, "retrireving entry")
	}

	return
}

// AppendEntries appends some entries to the shopping list
func (u *User) AppendEntries(names ...string) (err error) {
	// Ensures the user exists
	if _, err = GetUser("UID", u.UID); err != nil {
		return
	}

	// Prepares the statement
	stmt, err := db.Prepare(`INSERT INTO entries (uid, name) VALUES ($1, $2) ON CONFLICT DO NOTHING;`)
	defer stmt.Close()
	if err != nil {
		slog.Error("while preparing statement to append entries:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Inserts the entries
	for _, name := range names {
		if _, err = stmt.Exec(u.UID, name); err != nil {
			slog.Error("while appending shopping entry:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	return
}

// ToggleEntry toggles an entry's marked field
func (u *User) ToggleEntry(EID int) (err error) {
	// Makes sure the entry exists
	if _, err = u.GetEntry(EID); err != nil {
		return
	}

	// Updates it
	_, err = db.Exec(`UPDATE entries SET marked=(NOT marked) WHERE uid=$1 AND eid=$2;`, u.UID, EID)
	if err != nil {
		slog.Error("while toggling shopping entries:", "err", err)
		err = ERR_UNKNOWN
	}

	return
}

// ClearEntries drops all the marked entries
func (u *User) ClearShoppingList() (err error) {
	// Deletes the marked entries
	res, err := db.Exec(`DELETE FROM entries WHERE uid=$1 AND marked;`, u.UID)
	if err != nil {
		slog.Error("while clearing entries:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// Makes sure the user exists
		_, err = GetUser("UID", u.UID)
	}

	return
}

// EditEntry changes an entry's name
func (u *User) EditEntry(EID int, newName string) (err error) {
	// Gets the entry
	entry, err := u.GetEntry(EID)
	if err != nil {
		return
	}

	// Makes sure the new name is actually new
	if entry.Name == newName {
		return
	}

	// Makes sure the new name is not used
	var found int
	db.QueryRow(`SELECT 1 FROM entries WHERE uid=$1 AND name=$2;`, u.UID, newName).Scan(&found)
	if found > 0 {
		err = ERR_ENTRY_DUPLICATED
		return
	}

	// Change the name
	_, err = db.Exec(`UPDATE entries SET name=$3 WHERE uid=$1 AND eid=$2;`, u.UID, EID, newName)
	if err != nil {
		slog.Error("while editing entry:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}
