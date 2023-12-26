package mysql

import (
	"database/sql"
	"github.com/arithfi/arithfi-periphery/configs"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var (
	MYSQL *sql.DB
)

func init() {
	MYSQL = connectMysql()
}

func connectMysql() *sql.DB {
	db, err := sql.Open("mysql", configs.EnvMysqlURI())
	if err != nil {
		log.Fatal("Failed to connect to Mysql", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping Mysql: ", err)
	}
	log.Println("Connected to Mysql")

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
