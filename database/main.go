package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"strings"

	"cucinassistant/config"
)

var db *sql.DB

// Connect creates a connection to the database.
// It gets all the needed information from config.Runtime
func Connect() {
	// Connects to the database
	var err error
	db, err = sql.Open("postgres", config.Runtime.Database)
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
	for _, query := range strings.Split(schema, ";") {
		if strings.TrimSpace(query) != "" {
			if _, err := db.Exec(query + ";"); err != nil {
				slog.Error("while creating table:", "err", err)
				os.Exit(1)
			}
		}
	}
}

// handleNoRowsError is an utility function that does the following.
// If err is sql.ErrNoRows, checks if it happened because the user (with given
// UID) does not exist (and in this case it returns ERR_USER_UNKNOWN), or just because
// there are no rows (and in this case it returns ifExists).
// If err is not sql.ErrNoRows, it prints a log (with whileDoing) and returns
// ERR_UNKNOWN
func handleNoRowsError(err error, UID int, ifExist error, whileDoing string) error {
	if errors.Is(err, sql.ErrNoRows) {
		if _, err = GetUser("UID", UID); err == nil {
			return ifExist
		} else {
			return err
		}
	} else {
		slog.Error("while "+whileDoing+":", "err", err)
		return ERR_UNKNOWN
	}
}
