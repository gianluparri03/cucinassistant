package database

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/lib/pq"
)

var (
	// dateLocale is used to make tests work
	dateLocale = time.FixedZone("", 0)

	// defaultExpiration is used to match two items with the same name
	// and null expiration
	defaultExpiration = time.Date(3004, time.February, 5, 0, 0, 0, 0, dateLocale)
)

// Article is an item inside a storage section
// An article is identified by the name and the expiration.
// Both the expiration and the quantity can be null.
type Article struct {
	// SID is the Section ID
	SID int

	// AID is the Article ID
	AID int

	// Name is the article's name
	Name string

	// Expiration is the expiration date of the article.
	// It may be nil
	Expiration *time.Time

	// Quantity is the quantity of the article.
	// It may be nil
	Quantity *float32
}

// fixExpiration sets a nil expiration if it's the default
func (a *Article) fixExpiration() {
	if a != nil && *a.Expiration == defaultExpiration {
		a.Expiration = nil
	}
}

// FormatExpiration returns the expiration as a string
func (a Article) FormatExpiration() string {
	return a.Expiration.Format(time.DateOnly)
}

// IsExpired returns true if the article is expired
func (a Article) IsExpired() bool {
	return a.Expiration != nil && a.Expiration.Before(time.Now())
}

// Section is a named collection of articles
type Section struct {
	// Section is the Section ID
	SID int

	// Name is the name of the section
	Name string

	// Articles contains all the articles in this section
	Articles []Article
}

// StringArticle is a container for name, quantity,
// expiration and section as strings, used for inputs
type StringArticle struct {
	Section    string
	Name       string
	Quantity   string
	Expiration string
}

// Parse converts a StringArticle into an Article.
// Dates must be formatted like 2004-02-05 (time.DateOnly).
// If an empty date is given, the Article's expiration will
// be set to defaultExpiration, in order to make correct inserts
// in the database. When reading, instead, the Article.fixExpiration method
// will convert a defaultExpiration in a nil.
func (sa StringArticle) Parse() (Article, error) {
	var a Article

	// Parses the SID
	var err error
	if a.SID, err = strconv.Atoi(sa.Section); err != nil {
		return a, ERR_SECTION_NOT_FOUND
	}

	// Converts quantity to a float or nil
	if sa.Quantity == "" {
		a.Quantity = nil
	} else {
		if qty64, err := strconv.ParseFloat(sa.Quantity, 32); err == nil {
			qty32 := float32(qty64)
			a.Quantity = &qty32
		} else {
			return a, ERR_ARTICLE_QUANTITY_INVALID
		}
	}

	// Converts expiration to a time or nil
	if sa.Expiration == "" {
		a.Expiration = &defaultExpiration
	} else {
		if exp, err := time.ParseInLocation(time.DateOnly, sa.Expiration, dateLocale); err == nil {
			a.Expiration = &exp
		} else {
			return a, ERR_ARTICLE_EXPIRATION_INVALID
		}
	}

	a.Name = sa.Name
	return a, nil
}

// Storage is used to manage sections and articles
type Storage struct {
	uid int
}

// Storage returns the storage manager for the user
func (u User) Storage() Storage {
	return Storage{uid: u.UID}
}

// AddArticles adds some articles in multiple sections. If they are already
// present it will sum the quantities.
// If at least one of the two quantities is not given, the result
// will have the quantity unset.
func (s Storage) AddArticles(stringArticles ...StringArticle) error {
	var err error

	// Converts the string articles into articles
	articles := make([]Article, len(stringArticles))
	for i, sa := range stringArticles {
		if articles[i], err = sa.Parse(); err != nil {
			return err
		}
	}

	// Ensures that all the sections are owned by the user
	for _, a := range articles {
		if _, err := s.GetSection(a.SID); err != nil {
			return err
		}
	}

	// Prepares the statement
	stmt, err := db.Prepare(`INSERT INTO articles (sid, name, quantity, expiration) VALUES ($1, $2, $3, $4)
                             ON CONFLICT (sid, name, expiration) DO UPDATE set quantity = articles.quantity+excluded.quantity;`)
	defer stmt.Close()
	if err != nil {
		return ERR_UNKNOWN
	}

	// Inserts the entries
	for _, a := range articles {
		if _, err = stmt.Exec(a.SID, a.Name, a.Quantity, a.Expiration); err != nil {
			return ERR_UNKNOWN
		}
	}

	return nil
}

// DeleteArticle deletes an article
func (s Storage) DeleteArticle(AID int) error {
	// Makes sure the article exists and the user owns it
	if _, err := s.GetArticle(AID); err != nil {
		return err
	}

	// Deletes the article
	if _, err := db.Exec(`DELETE FROM articles WHERE aid=$1;`, AID); err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// DeleteSection tries to delete a section, with all the related articles
func (s Storage) DeleteSection(SID int) error {
	// Executes the query
	res, err := db.Exec(`DELETE FROM sections WHERE uid=$1 AND sid=$2;`, s.uid, SID)
	if err != nil {
		return ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the section (and the user) exist
		_, err = s.GetSection(SID)
		return err
	}

	return nil
}

// EditArticle tries to replace the article's name, quantity and/or expiration.
// It can also be used to move articles between sections.
func (s Storage) EditArticle(AID int, newData StringArticle) error {
	// Gets the current data
	article, err := s.GetArticle(AID)
	if err != nil {
		return err
	}

	// Parse the new data
	if parsed, err := newData.Parse(); err != nil {
		return err
	} else {
		// If nothing has changed, don't do anything
		if article.Name == parsed.Name &&
			article.SID == parsed.SID &&
			reflect.DeepEqual(article.Expiration, parsed.Expiration) &&
			reflect.DeepEqual(article.Quantity, parsed.Quantity) {
			return nil
		} else {
			article.SID = parsed.SID
			article.Name = parsed.Name
			article.Expiration = parsed.Expiration
			article.Quantity = parsed.Quantity
		}
	}

	// Checks if the section is owned by the user
	if _, err := s.GetSection(article.SID); err != nil {
		return err
	}

	// Checks if a similar article is already in storage
	var found int
	err = db.QueryRow(`SELECT 1 FROM articles
					   WHERE sid=$1 AND aid!=$2 AND name=$3 AND expiration=$4;`,
		article.SID, AID, article.Name, article.Expiration).Scan(&found)
	if found > 0 {
		return ERR_ARTICLE_DUPLICATED
	} else if err != nil {
		// If the error is sql.ErrNoRows, actually that's not an error, but
		// the desired output
		if !errors.Is(err, sql.ErrNoRows) {
			return ERR_UNKNOWN
		}
	}

	// Updates the article
	_, err = db.Exec(`UPDATE articles SET name=$3, expiration=$4, quantity=$5, sid=$2 WHERE aid=$1;`,
		AID, article.SID, article.Name, article.Expiration, article.Quantity)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// EditSection changes a section name
func (s Storage) EditSection(SID int, newName string) error {
	// Gets the section
	section, err := s.GetSection(SID)
	if err != nil {
		return err
	}

	// Makes sure the new name is actually new
	if section.Name == newName {
		return nil
	}

	// Makes sure the new name is not used
	var found int
	db.QueryRow(`SELECT 1 FROM sections WHERE uid=$1 AND name=$2;`, s.uid, newName).Scan(&found)
	if found > 0 {
		return ERR_SECTION_DUPLICATED
	}

	// Change the name
	_, err = db.Exec(`UPDATE sections SET name=$3 WHERE uid=$1 AND sid=$2;`, s.uid, SID, newName)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// GetArticle returns a specific article.
func (s Storage) GetArticle(AID int) (Article, error) {
	// Fetches the article
	var article Article
	err := db.QueryRow(`SELECT sid, aid, name, expiration, quantity FROM articles WHERE aid=$1;`, AID).
		Scan(&article.SID, &article.AID, &article.Name, &article.Expiration, &article.Quantity)

	if err != nil {
		return Article{}, handleNoRowsError(err, s.uid, ERR_ARTICLE_NOT_FOUND)
	} else {
		article.fixExpiration()
	}

	// Makes sure the section is owned by the user
	if _, err := s.GetSection(article.SID); err != nil {
		return Article{}, ERR_ARTICLE_NOT_FOUND
	}

	return article, nil
}

// GetArticles returns a section filled with its articles.
// If SID=0 it will return all the articles.
// The articles can be filtered by their name.
func (s Storage) GetArticles(SID int, nfilter string) (Section, error) {
	var section Section
	var sids []int

	if SID == 0 {
		// Gets all the user's sections' SID
		if sections, err := s.GetSections(); err == nil {
			for _, sec := range sections {
				sids = append(sids, sec.SID)
			}
		} else {
			return section, err
		}
	} else {
		// Ensures the section is owned by the user
		var err error
		if section, err = s.GetSection(SID); err == nil {
			sids = append(sids, section.SID)
		} else {
			return section, err
		}
	}

	// Runs the query
	rows, err := db.Query(`SELECT sid, aid, name, expiration, quantity
						   FROM articles WHERE sid = ANY($1) AND
						   name ILIKE CONCAT('%', $2::VARCHAR, '%')
						   ORDER BY expiration, aid;`, pq.Array(sids), nfilter)
	defer rows.Close()

	// Scans the articles
	if err != nil {
		return section, ERR_UNKNOWN
	} else {
		for rows.Next() {
			var a Article
			rows.Scan(&a.SID, &a.AID, &a.Name, &a.Expiration, &a.Quantity)
			a.fixExpiration()
			section.Articles = append(section.Articles, a)
		}
	}

	return section, nil
}

// GetNeighbours returns the previous and next articles' AIDS.
// If SID=0, it follows the global order, otherwise if follows the section
// order.
// It does not check for the article existence or ownership. Even if the SID
// is not the SID of the article's section, the results will be zero.
func (s Storage) GetNeighbours(SID int, AID int) (int, int) {
	// Prepares the sections to look into
	var sids []int
	if SID == 0 {
		if sections, err := s.GetSections(); err == nil {
			for _, section := range sections {
				sids = append(sids, section.SID)
			}
		}
	} else {
		sids = append(sids, SID)
	}

	// Runs the query
	var prev, next int
	db.QueryRow(`WITH ordered AS (SELECT aid,
						LAG(aid) OVER (ORDER BY expiration, aid) as prev,
						LEAD(aid) OVER (ORDER BY expiration, aid) as next
						FROM articles WHERE sid=ANY($1))
						SELECT COALESCE(prev, 0), COALESCE(next, 0)
						FROM ordered WHERE aid=$2;`,
		pq.Array(sids), AID,
	).Scan(&prev, &next)

	return prev, next
}

// GetSection returns a specific section, without the articles
func (s Storage) GetSection(SID int) (Section, error) {
	var section Section

	// Scans the section
	err := db.QueryRow(`SELECT sid, name FROM sections WHERE uid=$1 AND sid=$2;`, s.uid, SID).Scan(&section.SID, &section.Name)
	if err != nil {
		return section, handleNoRowsError(err, s.uid, ERR_SECTION_NOT_FOUND)
	}

	return section, nil
}

// GetSections returns all the sections created by an user.
// The articles are not fetched
func (s Storage) GetSections() ([]Section, error) {
	var sections []Section

	// Queries the sections
	rows, err := db.Query(`SELECT sid, name FROM sections WHERE uid=$1;`, s.uid)
	defer rows.Close()
	if err != nil {
		return sections, ERR_UNKNOWN
	}

	// Appends them to the list
	for rows.Next() {
		var s Section
		rows.Scan(&s.SID, &s.Name)
		sections = append(sections, s)
	}

	// If no sections have been found, makes sure the user exists
	if len(sections) == 0 {
		_, err = GetUser("UID", s.uid)
		return sections, err
	}

	return sections, nil
}

// NewSection tries to create a new section
func (s Storage) NewSection(name string) (Section, error) {
	// Ensures the user exists
	if _, err := GetUser("UID", s.uid); err != nil {
		return Section{}, err
	}

	// Checks if the name is used
	var found bool
	db.QueryRow(`SELECT 1 FROM sections WHERE uid=$1 AND name=$2;`, s.uid, name).Scan(&found)
	if found {
		return Section{}, ERR_SECTION_DUPLICATED
	}

	// Tries to save it in the database
	section := Section{Name: name}
	err := db.QueryRow(`INSERT INTO sections (uid, name) VALUES ($1, $2) RETURNING sid;`, s.uid, name).Scan(&section.SID)
	if err != nil {
		return Section{}, ERR_UNKNOWN
	}

	return section, nil
}
