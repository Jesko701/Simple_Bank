package main

import (
	"database/sql"
	"log"
	"os"
	"solo_simple-bank_tutorial/api"
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/util"

	_ "github.com/lib/pq"
)

func main() {
	//Load Configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	dbDriver := config.DBDriver
	if dbDriver == "" {
		dbDriver = os.Getenv("DB_Driver") // fallback to local environtment
		if dbDriver == "" {
			log.Fatal("Not Provide the DB Driver")
		}
	}

	dbSource := config.DBSource
	if dbSource == "" {
		dbSource = os.Getenv("DB_Source") // fallback to local environtment
		if dbSource == "" {
			log.Fatal("DB Source is not provided")
		}
	}

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	// Based on sqlc package
	store := sqlc.NewStore(conn)
	// Connect to router
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot create a server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
