package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chp0304/simplebank/api"
	db "github.com/chp0304/simplebank/db/sqlc"
	"github.com/chp0304/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	path := "."
	config, err := util.LoadConfig(path)
	fmt.Println(config)
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	server := api.NewServer(db.NewStore(conn))
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
