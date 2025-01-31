package database

import (
	"log/slog"
	"os"
	"testing"

	"cucinassistant/configs"
)

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelError)
	os.Chdir("..")

	// Loads the configuration
	configs.LoadAndParse()

	// Connects to the database and do the bootstrap
	Connect()
	Bootstrap()

	// Runs the actual tests
	m.Run()
}
