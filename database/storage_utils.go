package database

import (
	"strconv"
	"time"
)

// dateLocale is used to make tests work
var dateLocale = time.FixedZone("", 0)

// defaultExpiration is used to match two items with the same name
// and null expiration
var defaultExpiration = time.Date(2004, time.February, 5, 0, 0, 0, 0, dateLocale)

// StringArticle is a container for name, quantity
// and expiration as strings, used for inputs
type StringArticle struct {
	Name       string
	Quantity   string
	Expiration string
}

// Parse reads the data contained in the StringArticle
// and returns an Article.
// Dates must be formatted like 2004-02-05 (time.DateOnly).
// If an empty date is given, the Article's expiration will
// be set to defaultExpiration, in order to make correct inserts
// in the database. When reading, instead, the Article.fixExpiration method
// will convert a defaultExpiration in a nil.
func (sa StringArticle) Parse() (Article, error) {
	var a Article

	// Converts quantity to an int or nil
	if sa.Quantity == "" {
		a.Quantity = nil
	} else {
		if qty, err := strconv.Atoi(sa.Quantity); err == nil {
			a.Quantity = &qty
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

// fixExpiration sets a nil expiration if it's the default
func (a *Article) fixExpiration() {
	if a != nil && *a.Expiration == defaultExpiration {
		a.Expiration = nil
	}
}

// OrderedArticle is an article that can also
// store the AIDs of the previous and the next articles.
type OrderedArticle struct {
	Article
	Prev *int
	Next *int
}

// GetArticle returns a simple article
func (oa OrderedArticle) GetArticle() Article {
	return Article{AID: oa.AID, Name: oa.Name, Expiration: oa.Expiration, Quantity: oa.Quantity}
}
