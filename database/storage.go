package database

import (
	"database/sql"
	"log/slog"
	"time"
)

// Section is a named collection of articles
type Section struct {
	// Section is the Section ID
	SID int

	// Name is the name of the section
	Name string

	// Articles contains all the articles in this section
	Articles []Article
}

// GetSections returns all the sections created by an user.
// The articles are not fetched
func (u *User) GetSections() (sections []Section, err error) {
	// Queries the sections
	var rows *sql.Rows
	rows, err = db.Query(`SELECT sid, name FROM sections WHERE uid=$1;`, u.UID)
	if err != nil {
		slog.Error("while retrieving sections:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var s Section
		rows.Scan(&s.SID, &s.Name)
		sections = append(sections, s)
	}

	// If no sections have been found, makes sure the user exists
	if len(sections) == 0 {
		_, err = GetUser("UID", u.UID)
	}

	return
}

// NewSection tries to create a new section
func (u *User) NewSection(name string) (section Section, err error) {
	// Ensures the user exists
	if _, err = GetUser("UID", u.UID); err != nil {
		return
	}

	// Checks if the name is used
	var found bool
	db.QueryRow(`SELECT 1 FROM sections WHERE uid=$1 AND name=$2;`, u.UID, name).Scan(&found)
	if found {
		err = ERR_SECTION_DUPLICATED
		return
	}

	// Tries to save it in the database
	err = db.QueryRow(`INSERT INTO sections (uid, name) VALUES ($1, $2) RETURNING sid;`, u.UID, name).Scan(&section.SID)
	if err != nil {
		slog.Error("while creating section:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	section.Name = name
	return
}

// GetSection returns a specific section.
// If fetchArticles is true, it will also add the articles
// to the result. Filter
func (u *User) GetSection(SID int, fetchArticles bool) (section Section, err error) {
	// Scans the section
	err = db.QueryRow(`SELECT sid, name FROM sections WHERE uid=$1 AND sid=$2;`, u.UID, SID).Scan(&section.SID, &section.Name)
	if err != nil {
		err = handleNoRowsError(err, u.UID, ERR_SECTION_NOT_FOUND, "retrieving section")
		return
	}

	// if fetchArticles {
	// 	var rows *sql.Rows
	// 	rows, err = db.Query(`SELECT aid, name, expiration, quantity FROM articles WHERE sid=$1;`, SID)
	// 	if err != nil {
	// 		slog.Error("while retrieving articles:", "err", err)
	// 		err = ERR_UNKNOWN
	// 	} else {
	// 		defer rows.Close()
	// 		for rows.Next() {
	// 			var article Article
	// 			rows.Scan(article.Name, article.Expiration, article.Quantity)
	// 			section.Articles = append(section.Articles, &article)
	// 		}
	// 	}
	// }

	return
}

// EditSection changes a section name
func (u *User) EditSection(SID int, newName string) (err error) {
	// Gets the section
	var section Section
	if section, err = u.GetSection(SID, false); err != nil {
		return
	}

	// Makes sure the new name is actually new
	if section.Name == newName {
		return
	}

	// Makes sure the new name is not used
	var found int
	db.QueryRow(`SELECT 1 FROM sections WHERE uid=$1 AND name=$2;`, u.UID, newName).Scan(&found)
	if found > 0 {
		err = ERR_SECTION_DUPLICATED
		return
	}

	// Change the name
	_, err = db.Exec(`UPDATE sections SET name=$3 WHERE uid=$1 AND sid=$2;`, u.UID, SID, newName)
	if err != nil {
		slog.Error("while editing section:", "err", err)
		err = ERR_UNKNOWN
	}

	return
}

// DeleteSection tries to delete a section, with all the related articles
func (u *User) DeleteSection(SID int) (err error) {
	// Executes the query
	res, err := db.Exec(`DELETE FROM sections WHERE uid=$1 AND sid=$2;`, u.UID, SID)
	if err != nil {
		slog.Error("while deleting section:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the section (and the user) exist
		_, err = u.GetSection(SID, false)
	}

	return
}

// Article is an item inside a storage section
// An article is identified by the name and the expiration.
// Both the expiration and the quantity can be null.
type Article struct {
	// AID is the Article ID
	AID int

	// Name is the article's name
	Name string

	// Expiration is the expiration date of the article.
	// It may be null
	Expiration *time.Time

	// Quantity is the quantity of the article.
	// It may be null
	Quantity *int
}

// AddArticles adds some articles in a section. If they are already
// present it will sum the quantities.
// If at least one of the two quantities is not given, the result
// will not have the quantity set.
func (u *User) AddArticles(SID int, articles ...Article) (err error) {
	// Ensures the section exists
	if _, err = u.GetSection(SID, false); err != nil {
		return
	}

	// Prepares the statement
	stmt, err := db.Prepare(`INSERT INTO articles (sid, name, quantity, expiration) VALUES ($1, $2, $3, $4)
                             ON CONFLICT (sid, name, expiration) DO UPDATE set quantity = articles.quantity+excluded.quantity;`)
	defer stmt.Close()
	if err != nil {
		slog.Error("while preparing statement to add articles:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Inserts the entries
	for _, a := range articles {
		if _, err = stmt.Exec(SID, a.Name, a.Quantity, a.Expiration); err != nil {
			slog.Error("while adding article:", "err", err)
			err = ERR_UNKNOWN
		}
	}

	return
}
