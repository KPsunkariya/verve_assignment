package services

import (
	"log"
	"time"
)

// Deduplicate ID using Redis
func IsUniqueID(id string) bool {
	// Try to set the key with a TTL of 1 minute
	success, err := redisClient.SetNX(ctx, id, true, 1*time.Minute).Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return false // Assume duplicate if Redis fails
	}

	return success // True if the key was newly set, false if it already exists
}
