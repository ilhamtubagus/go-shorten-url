package server

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func ConnectRedisClient() *redis.Client {
	log.Println("connecting to Redis server")
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	_, err := rdbClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to Redis server")

	return rdbClient
}
