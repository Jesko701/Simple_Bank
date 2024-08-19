package sqlc

import (
	"database/sql"
	"log"
	"os"
	"solo_simple-bank_tutorial/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Error connect to config gile", err)
	}

	// Fallback to environment variables if the config file is missing or incomplete
	dbDriver := config.DBDriver
	if dbDriver == "" {
		dbDriver = os.Getenv("DB_DRIVER")
	}

	dbSource := config.DBSource
	if dbSource == "" {
		dbSource = os.Getenv("DB_SOURCE")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error when connecting to database", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
