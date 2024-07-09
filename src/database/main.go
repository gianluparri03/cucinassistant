package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
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
		log.Fatal("ERR: " + err.Error())
	}

	// Makes sure the connection is valid
	if err = DB.Ping(); err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func Bootstrap() {
	// Reads the schema file
	bytes, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}

	// Executes all the CREATE TABLEs
	queries := strings.Split(string(bytes), ";")
	for _, query := range queries {
		if strings.TrimSpace(query) != "" {
			if _, err := DB.Exec(query + ";"); err != nil {
				log.Fatal("ERR: " + err.Error())
			}
		}
	}
}
