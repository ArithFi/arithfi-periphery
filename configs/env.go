package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file")
	} else {
		log.Println("Loading .env file")
	}
}

func EnvMongoURI() string {
	// get url from .env
	url := os.Getenv("MONGODB_URI")
	return url
}

func EnvMysqlURI() string {
	url := os.Getenv("DSN")
	return url
}

func EnvRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	return addr
}

func EnvRedisPassword() string {
	password := os.Getenv("REDIS_PASSWORD")
	return password
}

func EnvBscScanAPI() string {
	url := os.Getenv("BSC_SCAN_API")
	return url
}
