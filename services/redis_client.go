package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedisClient(address string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     address, // Change to your Redis server address
		Password: "",      // No password by default
		DB:       0,       // Default DB
	})

	// Test the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully!")
}
