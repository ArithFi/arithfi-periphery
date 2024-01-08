package main

import (
	"database/sql"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", configs.EnvMysqlURI())
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

	log.Println("Successfully connected to Mysql!")

	// sql 查询 f_future_trading 表
	rows, err := db.Query("SELECT count(*) FROM f_future_trading")
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			log.Fatalf("failed to worker: %v", err)
		}
		fmt.Println(count)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("failed to iterate: %v", err)
	}
}
