package redis

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	REDIS *redis.Client
)

func init() {
	REDIS = connectCache()
}

func connectCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.EnvRedisURI(),
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis: ", err)
	}
	log.Println("Connected to Redis")
	return rdb
}
