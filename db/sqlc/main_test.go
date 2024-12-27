package db

import (
	"database/sql"
	"log"
	"os"
	"simple-bank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatalf("unable to load config file %e", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Panic("unable to open the database connection")
	}
	if err = testDb.Ping(); err != nil {
		log.Panic("unable to rach the database")
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
