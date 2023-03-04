package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // without calling any function of lib/pq, need to put _ in front of lib
)

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(DBDriver, DBSource) // in order to talk to specific db engine：go get github.com/lib/pq

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
