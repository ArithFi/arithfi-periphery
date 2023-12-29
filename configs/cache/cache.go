package cache

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	CACHE *redis.Client
)

func init() {
	CACHE = connectRedis()
}

func connectRedis() *redis.Client {
	opt, _ := redis.ParseURL(configs.EnvRedisURL())
	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis: ", err)
	}
	log.Println("Connected to Redis!")
	return rdb
}
