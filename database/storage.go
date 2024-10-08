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

// GetSection returns a specific section, without the articles
func (u *User) GetSection(SID int) (section Section, err error) {
	// Scans the section
	err = db.QueryRow(`SELECT sid, name FROM sections WHERE uid=$1 AND sid=$2;`, u.UID, SID).Scan(&section.SID, &section.Name)
	if err != nil {
		err = handleNoRowsError(err, u.UID, ERR_SECTION_NOT_FOUND, "retrieving section")
	}

	return
}

// EditSection changes a section name
func (u *User) EditSection(SID int, newName string) (err error) {
	// Gets the section
	var section Section
	if section, err = u.GetSection(SID); err != nil {
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
		_, err = u.GetSection(SID)
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
func (u *User) AddArticles(SID int, stringArticles ...StringArticle) (err error) {
	// Ensures the section exists
	if _, err = u.GetSection(SID); err != nil {
		return
	}

	// Converts the string articles into articles
	articles := make([]Article, len(stringArticles))
	for i, sa := range stringArticles {
		if articles[i], err = sa.Parse(); err != nil {
			return
		}
	}

	// Prepares the statement
	var stmt *sql.Stmt
	stmt, err = db.Prepare(`INSERT INTO articles (sid, name, quantity, expiration) VALUES ($1, $2, $3, $4)
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

// GetArticle returns a specific article
func (u *User) GetArticle(SID int, AID int) (article Article, err error) {
	// Makes sure the section is owned by the user
	if _, err = u.GetSection(SID); err != nil {
		return
	}

	// Fetches the article
	err = db.QueryRow(`SELECT aid, name, expiration, quantity FROM articles WHERE sid=$1 AND aid=$2;`, SID, AID).
		Scan(&article.AID, &article.Name, &article.Expiration, &article.Quantity)

	if err != nil {
		err = handleNoRowsError(err, u.UID, ERR_ARTICLE_NOT_FOUND, "retrieving article")
	} else {
		article.fixExpiration()
	}

	return
}

// GetArticles returns a specific section, filled with
// its articles. It is possible to specify a filter,
// that will filter the name.
func (u *User) GetArticles(SID int, filter string) (section Section, err error) {
	// Gets the empty section
	if section, err = u.GetSection(SID); err != nil {
		return
	}

	// Scans the articles
	var rows *sql.Rows
	rows, err = db.Query(`SELECT aid, name, expiration, quantity FROM articles
						  WHERE sid=$1 AND name ILIKE CONCAT('%', $2::VARCHAR, '%') ORDER BY expiration;`, SID, filter)
	if err != nil {
		slog.Error("while retrieving articles:", "err", err)
		err = ERR_UNKNOWN
	} else {
		defer rows.Close()
		for rows.Next() {
			var article Article
			rows.Scan(&article.AID, &article.Name, &article.Expiration, &article.Quantity)
			article.fixExpiration()
			section.Articles = append(section.Articles, article)
		}
	}

	return
}

// GetOrderedArticle returns a specific article, as an OrderedArticle
func (u *User) GetOrderedArticle(SID int, AID int) (article OrderedArticle, err error) {
	// Makes sure the section is owned by the user
	if _, err = u.GetSection(SID); err != nil {
		return
	}

	// Fetches the article and its neighbours aids
	err = db.QueryRow(`WITH ordered AS (SELECT aid, name, expiration, quantity,
						LAG(aid) OVER (PARTITION BY sid ORDER BY expiration) as prev,
						LEAD(aid) OVER (PARTITION BY sid ORDER BY expiration) as next
						FROM articles WHERE sid=$1) SELECT * FROM ordered WHERE aid=$2;`, SID, AID).
		Scan(&article.AID, &article.Name, &article.Expiration, &article.Quantity, &article.Prev, &article.Next)

	if err != nil {
		err = handleNoRowsError(err, u.UID, ERR_ARTICLE_NOT_FOUND, "retrieving ordered article")
	} else {
		article.fixExpiration()
	}

	return
}

// DeleteArticle deletes an article
func (u *User) DeleteArticle(SID int, AID int) (err error) {
	// Makes sure the section is owned by the user
	if _, err = u.GetSection(SID); err != nil {
		return
	}

	// Deletes the article
	res, err := db.Exec(`DELETE FROM articles WHERE sid=$1 AND aid=$2;`, SID, AID)
	if err != nil {
		slog.Error("while deleting article:", "err", err)
		err = ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the menu (and the user) exist
		_, err = u.GetArticle(SID, AID)
	}

	return
}
