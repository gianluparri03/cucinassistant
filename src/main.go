package main

import (
	"fmt"
	"log"
	"os"

	"cucinassistant/config"
	"cucinassistant/web"
)

func main() {
	// Ensures the config file is given
	if len(os.Args) < 2 {
		fmt.Println("Please provide a config file.")
		os.Exit(1)
	}

	log.Print("CucinAssistant (v" + config.Version + ")")

	// Parses the config
	log.Print("Reading configs... ")
	config.Read(os.Args[1])

	// Adds a listener for shutting down the server if it's on debug mode
	if config.Runtime.Debug {
		go func() {
			fmt.Scanln()
			log.Fatal("ERR: Keyboard interrupt")
		}()

		log.Print("Starting web server on " + config.Runtime.ServerAddress + " (press ENTER to stop)...")
	} else {
		log.Print("Starting web server on " + config.Runtime.ServerAddress + "...")
	}

	// Starts the server
	web.Start()
}
