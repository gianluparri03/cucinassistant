package database

import (
	"os"
	"testing"

	"cucinassistant/config"
)

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	// Loads the configuration
	config.Read("../" + os.Args[len(os.Args)-1])

	// Connects to the database and do the bootstrap
	Connect()
	Bootstrap()

	// Runs the actual tests
	m.Run()
}
