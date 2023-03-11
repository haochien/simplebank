package main

import (
	"database/sql"
	"log"

	"github.com/haochien/simplebank/api"
	db "github.com/haochien/simplebank/db/sqlc"
	"github.com/haochien/simplebank/util"
	_ "github.com/lib/pq" // important!! don't forget this import to connect go sql to specific engine
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
