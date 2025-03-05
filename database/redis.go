package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr, // Change if Redis is hosted elsewhere
		Password: "",   // No password by default
		DB:       0,    // Use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
}

// Add an allowed domain to Redis
func AddAllowedDomain(domain string) error {
	return RedisClient.SAdd(ctx, "allowed_domains", domain).Err()
}

// Remove an allowed domain from Redis
func RemoveAllowedDomain(domain string) error {
	return RedisClient.SRem(ctx, "allowed_domains", domain).Err()
}

// Get all allowed domains from Redis
func GetAllowedDomains() ([]string, error) {
	return RedisClient.SMembers(ctx, "allowed_domains").Result()
}
