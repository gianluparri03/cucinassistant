package database

import (
	"database/sql"
	"log/slog"
	"strings"
)

// Section is a named collection of articles
type Section struct {
	// Section is the Section ID
	SID int

	// Name is the name of the section
	Name string

	// Articles contains all the articles in this section
	Articles []struct{}
}

// GetSections returns all the sections created by an user.
// The articles are not fetched
func (u *User) GetSections() (sections []*Section, err error) {
	sections = []*Section{}

	// Queries the sections
	var rows *sql.Rows
	rows, err = DB.Query(`SELECT sid, name FROM sections WHERE uid=?;`, u.UID)
	if err != nil {
		slog.Error("while retrieving sections:", "err", err)
		return nil, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var s Section
		rows.Scan(&s.SID, &s.Name)
		sections = append(sections, &s)
	}

	// If no sections have been found, makes sure the user exists
	if len(sections) == 0 {
		if _, err = GetUser("UID", u.UID); err != nil {
			sections = nil
		}
	}

	return
}

// NewSection tries to create a new section
func (u *User) NewSection(name string) (section *Section, err error) {
	// Ensures the user exists
	_, err = GetUser("UID", u.UID)
	if err != nil {
		return
	}

	// Checks if the name is used
	var found bool
	DB.QueryRow(`SELECT 1 FROM sections WHERE uid=? AND name=?;`, name).Scan(&found)
	if found {
		err = ERR_SECTION_DUPLICATED
		return
	}

	// Tries to save it in the database
	_, err = DB.Exec(`INSERT INTO sections (uid, name) VALUES (?, ?);`, u.UID, name)
	if err != nil {
		slog.Error("while creating section:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Retrieves the SID
	section = &Section{Name: name}
	DB.QueryRow(`SELECT sid FROM sections WHERE uid=? AND name=?;`, u.UID, name).Scan(&section.SID)
	return
}

// GetSection returns a specific section, with its articles
// if fetchArticle is true
func (u *User) GetSection(SID int, fetchArticles bool) (section *Section, err error) {
	section = &Section{}

	// Scans the section
	err = DB.QueryRow(`SELECT sid, name FROM sections WHERE uid=? AND sid=?;`, u.UID, SID).Scan(&section.SID, &section.Name)
	if err != nil {
		section = nil

		if strings.HasSuffix(err.Error(), "no rows in result set") {
			// Makes sure the user exists if the section has not been found
			if _, err = GetUser("UID", u.UID); err == nil {
				err = ERR_SECTION_NOT_FOUND
			}
		} else {
			slog.Error("while retrieving section:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	return
}

// EditSection changes a section name
func (u *User) EditSection(SID int, newName string) (err error) {
	// Gets the section
	section, err := u.GetSection(SID, false)
	if err != nil {
		return
	}

	// Makes sure the new name is actually new
	if section.Name == newName {
		return
	}

	// Makes sure the new name is not used
	var found int
	DB.QueryRow(`SELECT 1 FROM sections WHERE uid=? AND name=?;`, u.UID, newName).Scan(&found)
	if found > 0 {
		err = ERR_SECTION_DUPLICATED
		return
	}

	// Change the name
	_, err = DB.Exec(`UPDATE sections SET name=? WHERE uid=? AND sid=?;`, newName, u.UID, SID)
	if err != nil {
		slog.Error("while editing section:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}

// DeleteSection tries to delete a section, with all the related articles
func (u *User) DeleteSection(SID int) (err error) {
	// Executes the query
	res, err := DB.Exec(`DELETE FROM sections WHERE uid=? AND sid=?;`, u.UID, SID)
	if err != nil {
		slog.Error("while deleting section:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the section (and the user) exist
		if _, err = GetUser("UID", u.UID); err == nil {
			err = ERR_SECTION_NOT_FOUND
		}
	}

	return
}
