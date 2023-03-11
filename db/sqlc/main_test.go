package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/haochien/simplebank/util"
	_ "github.com/lib/pq" // without calling any function of lib/pq, need to put _ in front of lib
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource) // in order to talk to specific db engine：go get github.com/lib/pq

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
