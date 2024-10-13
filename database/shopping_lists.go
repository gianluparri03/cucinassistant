package database

import (
	"database/sql"
	"log/slog"
)

// ShoppingList is used to manage the shopping list
type ShoppingList struct {
	uid int
}

// ShoppingList returns the shopping list manager for the user
func (u User) ShoppingList() ShoppingList {
	return ShoppingList{uid: u.UID}
}

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
func (sl ShoppingList) GetAll() (entries map[int]Entry, err error) {
	entries = map[int]Entry{}

	// Queries the entries
	var rows *sql.Rows
	rows, err = db.Query(`SELECT eid, name, marked FROM entries WHERE uid=$1;`, sl.uid)
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
		entries[e.EID] = e
	}

	// If no entries have been found, makes sure the user exists
	if len(entries) == 0 {
		_, err = GetUser("UID", sl.uid)
	}

	return
}

// GetOne returns a single entry of the shopping list
func (sl ShoppingList) GetOne(EID int) (e Entry, err error) {
	// Fetches them
	err = db.QueryRow(`SELECT eid, name, marked FROM entries WHERE uid=$1 AND eid=$2;`, sl.uid, EID).Scan(&e.EID, &e.Name, &e.Marked)
	if err != nil {
		err = handleNoRowsError(err, sl.uid, ERR_ENTRY_NOT_FOUND, "retrireving entry")
	}

	return
}

// Append appends some entries to the shopping list
func (sl ShoppingList) Append(names ...string) (err error) {
	// Ensures the user exists
	if _, err = GetUser("UID", sl.uid); err != nil {
		return
	}

	// Prepares the statement
	var stmt *sql.Stmt
	stmt, err = db.Prepare(`INSERT INTO entries (uid, name) VALUES ($1, $2) ON CONFLICT DO NOTHING;`)
	defer stmt.Close()
	if err != nil {
		slog.Error("while preparing statement to append entries:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Inserts the entries
	for _, name := range names {
		if _, err = stmt.Exec(sl.uid, name); err != nil {
			slog.Error("while appending shopping entry:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	return
}

// Toggle toggles an entry's marked field
func (sl ShoppingList) Toggle(EID int) (err error) {
	// Makes sure the entry exists
	if _, err = sl.GetOne(EID); err != nil {
		return
	}

	// Updates it
	_, err = db.Exec(`UPDATE entries SET marked=(NOT marked) WHERE uid=$1 AND eid=$2;`, sl.uid, EID)
	if err != nil {
		slog.Error("while toggling shopping entries:", "err", err)
		err = ERR_UNKNOWN
	}

	return
}

// Clear deletes all the marked entries
func (sl ShoppingList) Clear() (err error) {
	// Deletes the marked entries
	res, err := db.Exec(`DELETE FROM entries WHERE uid=$1 AND marked;`, sl.uid)
	if err != nil {
		slog.Error("while clearing entries:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// Makes sure the user exists
		_, err = GetUser("UID", sl.uid)
	}

	return
}

// Edit changes an entry's name
func (sl ShoppingList) Edit(EID int, newName string) (err error) {
	// Gets the entry
	entry, err := sl.GetOne(EID)
	if err != nil {
		return
	}

	// Makes sure the new name is actually new
	if entry.Name == newName {
		return
	}

	// Makes sure the new name is not used
	var found int
	db.QueryRow(`SELECT 1 FROM entries WHERE uid=$1 AND name=$2;`, sl.uid, newName).Scan(&found)
	if found > 0 {
		err = ERR_ENTRY_DUPLICATED
		return
	}

	// Change the name
	_, err = db.Exec(`UPDATE entries SET name=$3 WHERE uid=$1 AND eid=$2;`, sl.uid, EID, newName)
	if err != nil {
		slog.Error("while editing entry:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}
