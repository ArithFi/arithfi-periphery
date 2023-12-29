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
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.EnvRedisAddr(),
		Password: configs.EnvRedisPassword(),
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis: ", err)
	}
	log.Println("Connected to Redis!")
	return rdb
}
