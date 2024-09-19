package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"strings"

	"cucinassistant/config"
)

var DB *sql.DB

// Connect creates a connection to the database.
// It gets all the needed information from config.Runtime
func Connect() {
	// Connects to the database
	var err error
	DB, err = sql.Open("postgres", config.Runtime.Database)
	if err != nil {
		slog.Error("while connecting to the db:", "err", err)
		os.Exit(1)
	}

	// Makes sure the connection is valid
	if err = DB.Ping(); err != nil {
		slog.Error("while pinging the db:", "err", err)
		err = DB.Ping()
		os.Exit(1)
	}
}

// Bootstrap makes sure that the database schema is ready
func Bootstrap() {
	for _, query := range strings.Split(schema, ";") {
		if strings.TrimSpace(query) != "" {
			if _, err := DB.Exec(query + ";"); err != nil {
				slog.Error("while creating table:", "err", err)
				os.Exit(1)
			}
		}
	}
}
