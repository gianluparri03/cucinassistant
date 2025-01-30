package database

import (
	"os"
	"testing"

	"cucinassistant/configs"
)

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	os.Chdir("..")

	// Loads the configuration
	configs.LoadAndParse()

	// Connects to the database and do the bootstrap
	Connect()
	Bootstrap()

	// Runs the actual tests
	m.Run()
}
