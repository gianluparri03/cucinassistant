package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/web"
)

// main runs everything
func main() {
	// Prints a welcome text
	title := fmt.Sprintf("CucinAssistant [v%s]", configs.Version)
	slog.Warn(title)
	slog.Warn(strings.Repeat("=", len(title)))

	// Loads the configs
	slog.Warn("Loading configs...")
	configs.LoadAndParse()

	// Connects to the database
	slog.Warn("Connecting to the database...")
	database.Connect()

	// Checks the schema
	slog.Warn("Checking schema...")
	database.Bootstrap()

	// Adds a listener for shutting down the server if it's on debug mode
	slog.Warn("Starting web server...")
	if configs.Debug {
		go func() {
			fmt.Scanln()
			os.Exit(0)
		}()

		slog.Warn("[Press ENTER to stop]")
	}

	// Starts the server
	slog.Warn("Running on " + configs.BaseURL)
	web.Start()
}
