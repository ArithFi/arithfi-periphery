package configs

import (
	"os"
)

func EnvMongoURI() string {
	// get url from .env
	url := os.Getenv("MONGODB_URI")
	return url
}
