// Cache Helpers
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/go-redis/redis/v8"
)

// RedisHelper wraps the Redis client and context for easier access
type RedisHelper struct {
	Client *redis.Client
	Ctx    context.Context
}

// NewRedisHelper initializes and returns a RedisHelper instance
func NewRedisHelper(cfg config.RedisConfigs) (*RedisHelper, error) {
	ctx := context.Background()

	address := "localhost:6379"

	client := redis.NewClient(&redis.Options{
		Addr: address,
		DB:   cfg.RedisDB, // Select the Redis database
	})
	// Test the connection with a timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if ping, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Error connecting to Redis at %s: %v", address, err)
		return nil, fmt.Errorf("redis connection error: %w", err)
	} else {
		log.Printf("Successfully connected to Redis at %s: %s", address, ping)
	}

	return &RedisHelper{
		Client: client,
		Ctx:    ctx,
	}, nil
}

// CacheGet retrieves the cached data for a given key.
func CacheGet(ctx context.Context, redisClient *redis.Client, key string, dest interface{}) error {
	cachedData, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil { // Cache miss
		return err
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cachedData), dest)
}

// CacheSet sets the cache with a TTL for the given key.
func CacheSet(ctx context.Context, redisClient *redis.Client, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, jsonData, ttl).Err()
}

func DeleteFeedEntry(ctx context.Context, client *redis.Client, key string) error {
	result, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		log.Printf("No entry found for key: %s", key)
	} else {
		log.Printf("Deleted key: %s", key)
	}

	return nil
}
