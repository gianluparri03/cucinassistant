package database

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"
	"testing"

	"cucinassistant/configs"
)

// testingDBName is the name of the database where the tests will be run.
// It is set to the database name (see configs.Database) with a _test suffix.
// The database is created and dropped at every test.
var testingDBName string

// testingOrigConn is used to hold the connection to the default database
// when connecting to the testing one.
var testingOrigConn *sql.DB

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelError)

	os.Chdir("..")
	configs.LoadAndParse()

	openTestDB()
	defer closeTestDB()

	Bootstrap()
	m.Run()
}

// openTestDB creates a testing database
func openTestDB() {
	var dbName string

	// Builds the testing db's name
	for _, pair := range strings.Split(configs.Database, " ") {
		if strings.HasPrefix(pair, "dbname=") {
			dbName = pair[len("dbname="):]
			break
		}
	}

	// Connects to the original database
	Connect()
	testingOrigConn = db

	// Creates the testing db
	testingDBName = dbName + "_test"
	_, err := db.Exec("CREATE DATABASE " + testingDBName + ";")
	if err != nil {
		slog.Error("while creating testing db:", "err", err)
		os.Exit(1)
	}

	// Connects to the testing database
	configs.Database = strings.ReplaceAll(configs.Database,
		"dbname="+dbName,
		"dbname="+testingDBName)
	Connect()
}

// closeTestDB drops the testing database
func closeTestDB() {
	db.Close()
	db = testingOrigConn

	// Drops the testing database
	_, err := db.Exec("DROP DATABASE " + testingDBName + ";")
	if err != nil {
		slog.Error("while dropping testing db:", "err", err)
		os.Exit(1)
	}
}
