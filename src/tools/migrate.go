package main

import (
	"fmt"
	"os"
	"log/slog"

	"cucinassistant/configs"
	"cucinassistant/database"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelError)

	// Prints a welcome text
	fmt.Print(`CucinAssistant Migration Tool
=============================
This tool MUST BE used once and only once, as it will make changes the database schema.
Type BOOTSTRAP to create it from scratch, or UPGRADE to upgrade it from the previous version
(that is, from version 8 (Banana) to version 9 (Maracuja)). To upgrade from previous versions
please perform each upgrade version by version.
`)

	// Asks a confirm
	var command string
	fmt.Scanln(&command)
	if command != "BOOTSTRAP" && command != "UPGRADE" {
		os.Exit(1)
	} else {
		fmt.Printf("Running %s.\n", command)
	}

	// Initializes all the modules
	configs.LoadAndParse()
	db := database.Connect()

	// Runs the command
	if command == "BOOTSTRAP" {
		database.Bootstrap()
	} else {
		// Adds the schema version
		db.Exec(`CREATE TABLE ca_version (id INT NOT NULL);`)
		db.Exec(`INSERT INTO ca_version VALUES ($1);`, configs.VersionCode)
	}

	fmt.Println("Done.")
}
