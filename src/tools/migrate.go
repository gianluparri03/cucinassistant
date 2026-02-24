package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"cucinassistant/configs"
	"cucinassistant/database"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelError)

	// Prints a welcome text
	fmt.Println("CucinAssistant Migration Tool")
	fmt.Println("=============================")

	// Initializes all the modules
	configs.LoadAndParse()
	db := database.Connect()
	fmt.Println("Connected to the database.")

	// Checks the database version, or if it has been initialized
	var version, initialized int
	db.QueryRow(`SELECT id FROM ca_version LIMIT 1;`).Scan(&version)
	if version > 0 {
		fmt.Printf("Your database appears to be from version %d.\n", version)
	} else {
		db.QueryRow(`SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' LIMIT 1;`).Scan(&initialized)
		if initialized == 1 {
			fmt.Println("Your database is from an unknown version.")
			fmt.Println("Please insert manually your schema version:")
			fmt.Scanf("%d", &version)
		} else {
			fmt.Println("Your database is empty.")
		}
	}

	// Tell the user the action that will be performed
	fmt.Println()
	if version == configs.VersionCode {
		fmt.Println("Your database schema is already updated to the latest version.")
		os.Exit(0)
	} else if initialized != 1 {
		fmt.Println("Please type CONFIRM to set up the database schema for the first time")
	} else if version >= 8 {
		fmt.Printf("Please type CONFIRM to update your database from version %d to version %d\n", version, configs.VersionCode)
	} else {
		fmt.Println("Unsupported version code. Please refer to the online documentation.")
		os.Exit(1)
	}

	// Asks a confirm
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "CONFIRM" {
		os.Exit(1)
	}

	// Upgrade (or set up) the schema
	fmt.Println()
	if initialized == 1 {
		switch version { // remember to add fallthrough
		case 8:
			update8to9(db)
		}
	} else {
		fmt.Println("Setting up the schema...")
		database.Bootstrap()
	}

	fmt.Println("Done.")
}

func update8to9(db *sql.DB) {
	fmt.Println("Upgrading from version 8 to version 9...")

	// Adds the schema version
	db.Exec(`CREATE TABLE ca_version (id INT NOT NULL);`)
	db.Exec(`INSERT INTO ca_version VALUES ($1);`, configs.VersionCode)

	// Adds the recipe tags
	db.Exec(`CREATE TABLE tags (name VARCHAR NOT NULL, rid INT NOT NULL, PRIMARY KEY (name, rid), FOREIGN KEY (rid) REFERENCES recipes (rid) ON DELETE CASCADE);`)
	db.Exec(`CREATE INDEX tags_name ON tags (name);`)
}
