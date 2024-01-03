package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	opt, _ := redis.ParseURL(configs.EnvRedisURL())
	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis: ", err)
	}
	log.Println("Connected to Redis!")
}
