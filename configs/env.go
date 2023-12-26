package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loading .env file")
}

func EnvMongoURI() string {
	// get url from .env
	url := os.Getenv("MONGODB_URI")
	return url
}

func EnvRedisURI() string {
	// get url from .env
	url := os.Getenv("REDIS_URI")
	return url
}

func EnvMysqlURI() string {
	url := os.Getenv("DSN")
	return url
}
