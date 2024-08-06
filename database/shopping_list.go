package database

import (
	"database/sql"
	"log/slog"
	"strings"
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

// ShoppingList contains the pairs EID: Entry
type ShoppingList map[int]Entry

// GetShoppingList returns an user's shopping list
func (u *User) GetShoppingList() (sl ShoppingList, err error) {
	sl = ShoppingList{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = DB.Query(`SELECT eid, name, marked FROM shopping_entries WHERE uid = ?;`, u.UID)
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
		_, err = GetUser(u.UID)
	}

	return
}

// GetEntry returns a single entry of the shopping list
func (u *User) GetEntry(EID int) (e Entry, err error) {
	err = DB.QueryRow(`SELECT eid, name, marked FROM shopping_entries WHERE uid = ? AND eid = ?;`, u.UID, EID).Scan(&e.EID, &e.Name, &e.Marked)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			err = ERR_ENTRY_NOT_FOUND
		} else {
			slog.Error("while retrieving shopping entry:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	return
}

// AppendEntries appends some entries to the shopping list
func (u *User) AppendEntries(names ...string) (err error) {
	// Ensures the user exists
	_, err = GetUser(u.UID)
	if err != nil {
		return
	}

	// Prepares the statement
	stmt, _ := DB.Prepare(`INSERT IGNORE INTO shopping_entries (uid, eid, name)
	                       SELECT ?, IFNULL(MAX(eid), 0) + 1, ? FROM shopping_entries WHERE uid=?;`)

	// Inserts the entries
	for _, name := range names {
		_, err_ := stmt.Exec(u.UID, name, u.UID)
		if err_ != nil {
			slog.Error("while appending shopping entries:", "err", err_)
			err = ERR_UNKNOWN
		}
	}

	return
}

// ToggleEntry toggles an entry's marked field
func (u *User) ToggleEntry(EID int) (err error) {
	res, err := DB.Exec(`UPDATE shopping_entries SET marked = !marked WHERE uid = ? AND eid = ?;`, u.UID, EID)
	if err != nil {
		slog.Error("while toggling shopping entries:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// Makes sure the entry has been toggled
		err = ERR_ENTRY_NOT_FOUND
	}

	return
}

// ClearEntries drops all the marked entries
func (u *User) ClearEntries() (err error) {
	res, err := DB.Exec(`DELETE FROM shopping_entries WHERE uid = ? AND marked;`, u.UID)
	if err != nil {
		slog.Error("while clearing entries:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// Makes sure the user exists
		_, err = GetUser(u.UID)
	}

	return
}

// EditEntry changes an entry's name
func (u *User) EditEntry(EID int, newName string) error {
	// Gets the entry
	entry, err := u.GetEntry(EID)
	if err != nil {
		return err
	}

	// Makes sure the new name is actually new
	if entry.Name == newName {
		return nil
	}

	// Makes sure the new name is not used
	var found int
	DB.QueryRow(`SELECT 1 FROM shopping_entries WHERE uid = ? AND name = ?;`, u.UID, newName).Scan(&found)
	if found > 0 {
		return ERR_ENTRY_DUPLICATED
	}

	// Change the name
	_, err = DB.Exec(`UPDATE shopping_entries SET name = ? WHERE uid = ? AND eid = ?;`, newName, u.UID, EID)
	if err != nil {
		slog.Error("while editing entry:", "err", err)
		return ERR_UNKNOWN
	}

	return nil
}
