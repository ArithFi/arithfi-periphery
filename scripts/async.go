package main

import (
	_ "github.com/arithfi/arithfi-periphery/configs"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"log"
)

// main function to async mysql => mongoDB
func main() {
	log.Println("Async Works Start!")
	err := mysql.MYSQL.Ping()
	if err != nil {
		return
	}
	// sql 查询 f_future_trading 表
	rows, err := mysql.MYSQL.Query("select ID from f_future_trading limit 10")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ID int
		if err := rows.Scan(&ID); err != nil {
			log.Fatalf("failed to scan: %v", err)
		}
		log.Println(ID)
	}
}
