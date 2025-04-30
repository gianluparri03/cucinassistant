package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"
	_ "github.com/lib/pq"
	_ "embed"

	"cucinassistant/configs"
)

var db *sql.DB

//go:embed schema.sql
var schema string

// handleNoRowsError is an utility function that does the following.
// If err is sql.ErrNoRows, checks if it happened because the user (with given
// UID) does not exist (and in this case it returns ERR_USER_UNKNOWN), or
// just because there are no rows (and in this case it returns ifExists).
// Otherwise returns ERR_UNKNOWN.
func handleNoRowsError(err error, UID int, ifExist error) error {
	if errors.Is(err, sql.ErrNoRows) {
		if _, err = GetUser("UID", UID); err == nil {
			return ifExist
		} else {
			return err
		}
	} else {
		return ERR_UNKNOWN
	}
}

// Connect creates a connection to the database.
func Connect() {
	// Connects to the database
	var err error
	db, err = sql.Open("postgres", configs.Database)
	if err != nil {
		slog.Error("while connecting to the db:", "err", err)
		os.Exit(1)
	}

	// Makes sure the connection is valid
	if err = db.Ping(); err != nil {
		slog.Error("while pinging the db:", "err", err)
		os.Exit(1)
	}
}

// Bootstrap makes sure that the database schema is ready
func Bootstrap() {
	// Splits it and executes it
	for _, query := range strings.Split(string(schema), ";") {
		if strings.TrimSpace(query) != "" {
			if _, err := db.Exec(query + ";"); err != nil {
				slog.Error("while creating table:", "err", err)
				os.Exit(1)
			}
		}
	}
}

// Stats is a report of the current database population
type Stats struct {
	UsersNumber    int
	MenusNumber    int
	SectionsNumber int
	ArticlesNumber int
	EntriesNumber  int
	RecipesNumber  int
}

// GetStats returns a Stats instance
func GetStats() (s Stats) {
	// Counts the records
	db.QueryRow(`SELECT COUNT(*) FROM ca_users;`).Scan(&s.UsersNumber)
	db.QueryRow(`SELECT COUNT(*) FROM menus;`).Scan(&s.MenusNumber)
	db.QueryRow(`SELECT COUNT(*) FROM sections;`).Scan(&s.SectionsNumber)
	db.QueryRow(`SELECT COUNT(*) FROM articles;`).Scan(&s.ArticlesNumber)
	db.QueryRow(`SELECT COUNT(*) FROM entries;`).Scan(&s.EntriesNumber)
	db.QueryRow(`SELECT COUNT(*) FROM recipes;`).Scan(&s.RecipesNumber)
	return
}
