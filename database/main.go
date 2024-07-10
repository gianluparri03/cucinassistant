package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log/slog"
	"os"
	"strings"

	"cucinassistant/config"
)

var DB *sql.DB

func Connect() {
	// Prepares the configuration
	credentials := mysql.Config{
		Addr:                 config.Runtime.Database.Address,
		User:                 config.Runtime.Database.Username,
		Passwd:               config.Runtime.Database.Password,
		DBName:               config.Runtime.Database.Database,
		AllowNativePasswords: true,
	}

	// DBects to the database
	var err error
	DB, err = sql.Open("mysql", credentials.FormatDSN())
	if err != nil {
		slog.Error("while connecting:", "err", err)
		os.Exit(1)
	}

	// Makes sure the connection is valid
	if err = DB.Ping(); err != nil {
		slog.Error("while pinging:", "err", err)
		err = DB.Ping()
		os.Exit(1)
	}
}

func Bootstrap(scriptPath string) {
	// Reads the schema file
	bytes, err := os.ReadFile(scriptPath)
	if err != nil {
		slog.Error("while reading schema:", "err", err)
		os.Exit(1)
	}

	// Executes all the CREATE TABLEs
	queries := strings.Split(string(bytes), ";")
	for _, query := range queries {
		if strings.TrimSpace(query) != "" {
			if _, err := DB.Exec(query + ";"); err != nil {
				slog.Error("while creating table:", "err", err)
				os.Exit(1)
			}
		}
	}
}
