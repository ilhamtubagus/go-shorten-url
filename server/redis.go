package server

import (
	"context"
	"fmt"
	"github.com/ilhamtubagus/go-shorten-url/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func ConnectRedisClient(config config.RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Println("connecting to Redis server", addr)

	rdbClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       0, // use default DB
	})

	_, err := rdbClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to Redis server")

	return rdbClient
}
