package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(opts *redis.Options) *RedisCache {
	client := redis.NewClient(opts)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to verify the connection to Redis")
	}

	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) CacheData(ctx context.Context, key, data string, duration time.Duration) (string, error) {
	return r.client.Set(ctx, key, data, duration).Result()
}

func (r *RedisCache) CheckCache(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) GracefulShutdown() {
	fmt.Println("Closing Redis Client")
	r.client.Close()
	fmt.Println("Closed Redis Client")
}
