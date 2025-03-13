package cache

import (
	"context"
	"encoding/json"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type CourseCache struct {
	redisClient *redis.Client
}

func NewCourseCache(redisClient *redis.Client) *CourseCache {
	return &CourseCache{
		redisClient: redisClient,
	}
}

func (c *CourseCache) Get(ctx context.Context, key string) ([]*entity.Course, error) {
	data, err := c.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var courses []*entity.Course
	err = json.Unmarshal([]byte(data), &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (c *CourseCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
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
