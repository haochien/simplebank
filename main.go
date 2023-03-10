package main

import (
	"database/sql"
	"log"

	"github.com/haochien/simplebank/api"
	db "github.com/haochien/simplebank/db/sqlc"
	_ "github.com/lib/pq" // important!! don't forget this import to connect go sql to specific engine
)

const (
	DBDriver      = "postgres"
	DBSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
