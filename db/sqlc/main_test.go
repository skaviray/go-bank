package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:admin@localhost:5433/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	// log.Println(testDb)
	if err != nil {
		log.Panic("unable to open the database connection")
	}
	if err = testDb.Ping(); err != nil {
		log.Panic("unable to rach the database")
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
