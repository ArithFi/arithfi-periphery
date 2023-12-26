package main

import (
	_ "github.com/arithfi/arithfi-periphery/configs"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"log"
)

// 这是一个脚本，用于同步ArithFi的数据库，并且建立新的表结构

// 不需要在运行时输入参数

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
