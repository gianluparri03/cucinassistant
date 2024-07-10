package database

import (
	"os"
	"testing"

	"cucinassistant/config"
)

func TestMain(m *testing.M) {
	// Loads the configuration
    config.Read(os.Args[len(os.Args)-1])

	// Connects to the database
	// and creates the missing tables
	Connect()
	Bootstrap("schema.sql")

	// Runs the actual tests
	m.Run()
}
