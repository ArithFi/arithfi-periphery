package configs

import (
	"os"
)

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
