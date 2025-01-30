package configs

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// LoadAndParse loads and parses the config files
func LoadAndParse() {
	load()
	parse()

	// Sets the logger level according to what has been read
	if Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}

// load loads the configs from
// * /configs/.env.E.local
// * /configs/.env.E
// * /configs/.env
// where E is read from CA_ENV (should be 'development', 'testing' or 'production').
// The former has the most priority.
func load() {
	env := os.Getenv("CA_ENV")
	godotenv.Load("configs/.env." + env + ".local")
	godotenv.Load("configs/.env." + env)
	godotenv.Load("configs/.env")

	slog.Warn("Environment", "type", env)
}

// parse ensures that all the required fields are loaded
func parse() {
	Test = parseBool("CA_TEST", false)
	Debug = parseBool("CA_DEBUG", false)
	BaseURL = parseString("CA_BASEURL", !Test)
	Port = parseString("CA_PORT", !Test)
	SessionSecret = parseString("CA_SESSIONSECRET", !Test)
	Database = parseString("CA_DATABASE", true)
	EmailEnabled = parseBool("CA_EMAIL_ENABLED", !Test)
	EmailSender = parseString("CA_EMAIL_SENDER", EmailEnabled)
	EmailServer = parseString("CA_EMAIL_SERVER", EmailEnabled)
	EmailPort = parseString("CA_EMAIL_PORT", EmailEnabled)
	EmailLogin = parseString("CA_EMAIL_LOGIN", EmailEnabled)
	EmailPassword = parseString("CA_EMAIL_PASSWORD", EmailEnabled)
}

// parseString reads a string from the environment variables, and
// shows an error if it's required and not set.
func parseString(env string, required bool) string {
	value, found := os.LookupEnv(env)
	if !found && required {
		slog.Error("required config field not set:", "field", env)
		os.Exit(1)
	}

	return value
}

// parseBool reads a bool from the environment variables, and
// shows an error if it's required and not set, or not a bool.
func parseBool(env string, required bool) bool {
	value := parseString(env, required)
	if value == "true" || value == "1" {
		return true
	} else if value == "false" || value == "0" || value == "" {
		return false
	} else {
		slog.Error("unknown config value:", "field", env)
		os.Exit(1)
	}

	return false
}
