package test

import (
	"os"
	"testing"

	"cucinassistant/config"
	"cucinassistant/database"
)

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	// Loads the configuration
	config.Read(os.Args[len(os.Args)-1])

	// Connects to the database
	// and creates the missing tables
	database.Connect()
	database.Bootstrap("../schema.sql")

	// Runs the actual tests
	m.Run()
}