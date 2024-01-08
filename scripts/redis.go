package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.EnvRedisAddr(),
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis: ", err)
	}
	log.Println("Connected to Redis!")
}
