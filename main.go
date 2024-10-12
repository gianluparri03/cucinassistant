package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web"
)

// main runs everything
func main() {
	// Ensures the config file is given
	if len(os.Args) < 2 {
		fmt.Println("Please provide a config file.")
		os.Exit(1)
	}

	// Prints a welcome text
	title := fmt.Sprintf("CucinAssistant (version: %s)", config.Version)
	slog.Warn(title)
	slog.Warn(strings.Repeat("=", len(title)))

	// Parses the config
	slog.Warn("Reading configs...")
	config.Read(os.Args[1])

	// Connects to the database
	slog.Warn("Connecting to the database...")
	database.Connect()

	// Checks the schema
	slog.Warn("Checking schema...")
	database.Bootstrap()

	// Adds a listener for shutting down the server if it's on debug mode
	slog.Warn("Starting web server...")
	if config.Runtime.Debugging {
		go func() {
			fmt.Scanln()
			os.Exit(0)
		}()

		slog.Warn("[Press ENTER to stop]")
	}

	// Starts the server
	slog.Warn("Running on http://localhost:" + config.Runtime.Port + "/")
	web.Start()
}
