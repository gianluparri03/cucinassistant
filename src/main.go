package main

import (
	"fmt"
	"log/slog"
	"os"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web"
)

func main() {
	// Ensures the config file is given
	if len(os.Args) < 2 {
		fmt.Println("Please provide a config file.")
		os.Exit(1)
	}

	slog.Info("CucinAssistant (v" + config.Version + ")")
	slog.Info("=======================")

	// Parses the config
	slog.Info("Reading configs...")
	config.Read(os.Args[1])

	// Connects to the database
	slog.Info("Connecting to the database...")
	database.Connect()

	// Creates missing tables
	slog.Info("Checking tables...")
	database.Bootstrap()

	// Adds a listener for shutting down the server if it's on debug mode
	if config.Runtime.Debug {
		go func() {
			fmt.Scanln()
			slog.Error("Keyboard interrupt")
			os.Exit(0)
		}()

		slog.Info("Starting web server on " + config.Runtime.ServerAddress + " (press ENTER to stop)...")
	} else {
		slog.Info("Starting web server on " + config.Runtime.ServerAddress + "...")
	}

	// Starts the server
	web.Start()
}
