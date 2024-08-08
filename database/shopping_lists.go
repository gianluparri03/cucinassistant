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

// GetShoppingList returns an user's shopping list
func (u *User) GetShoppingList() (sl map[int]*Entry, err error) {
	sl = map[int]*Entry{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = DB.Query(`SELECT eid, name, marked FROM shopping_entries WHERE uid = ?;`, u.UID)
	if err != nil {
		slog.Error("while retrieving shopping list:", "err", err)
		return nil, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var e Entry
		rows.Scan(&e.EID, &e.Name, &e.Marked)
		sl[e.EID] = &e
	}

	// If no entries have been found, makes sure the user exists
	if len(sl) == 0 {
		if _, err = GetUser("UID", u.UID); err != nil {
			sl = nil
		}
	}

	return
}

// GetEntry returns a single entry of the shopping list
func (u *User) GetEntry(EID int) (e *Entry, err error) {
	e = &Entry{}

	err = DB.QueryRow(`SELECT eid, name, marked FROM shopping_entries WHERE uid = ? AND eid = ?;`, u.UID, EID).Scan(&e.EID, &e.Name, &e.Marked)
	if err != nil {
		e = nil

		if strings.HasSuffix(err.Error(), "no rows in result set") {
			// Makes sure the user exists
			if _, err = GetUser("UID", u.UID); err == nil {
				err = ERR_ENTRY_NOT_FOUND
			}
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
	_, err = GetUser("UID", u.UID)
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
	// Gets the entry (to make sure it exists)
	if _, err = u.GetEntry(EID); err != nil {
		return
	}

	res, err := DB.Exec(`UPDATE shopping_entries SET marked = !marked WHERE uid = ? AND eid = ?;`, u.UID, EID)
	if err != nil {
		slog.Error("while toggling shopping entries:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		err = ERR_UNKNOWN
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
		if _, err = GetUser("UID", u.UID); err == nil {
			err = ERR_ENTRY_NOT_FOUND
		}
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
	DB.QueryRow(`SELECT 1 FROM shopping_entries WHERE uid = ? AND name = ?;`, u.UID, newName).Scan(&found)
	if found > 0 {
		err = ERR_ENTRY_DUPLICATED
		return
	}

	// Change the name
	_, err = DB.Exec(`UPDATE shopping_entries SET name = ? WHERE uid = ? AND eid = ?;`, newName, u.UID, EID)
	if err != nil {
		slog.Error("while editing entry:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}
