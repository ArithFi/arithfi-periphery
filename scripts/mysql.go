package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer func(db *sql.DB) {
		log.Println("Closing database connection")
		err := db.Close()
		if err != nil {
			log.Printf("failed to close: %v\n", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	log.Println("Successfully connected to PlanetScale!")
}
