package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

// Read reads the config from the path specified.
// The result is stored in the Runtime global variable.
func Read(path string) {
	var data []byte
	var err error

	// Reads the file
	if data, err = os.ReadFile(path); err != nil {
		slog.Error("while opening file:", "err", err)
		os.Exit(1)
	}

	// Parses it
	if err = yaml.Unmarshal(data, &Runtime); err != nil {
		slog.Error("while unmarshaling:", "err", err)
		os.Exit(1)
	}

	// Sets the logger level according to what has been read
	if Runtime.Mode == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else if Runtime.Mode == "test" {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
