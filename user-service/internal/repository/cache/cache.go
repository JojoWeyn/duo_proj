package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(redisClient *redis.Client) *RedisCache {
	return &RedisCache{
		redisClient: redisClient,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.redisClient.Set(ctx, key, data, ttl).Err()
	if err != nil {
		log.Printf("Error setting cache data: %v", err)
		return err
	}

	return nil
}
