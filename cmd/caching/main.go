package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	rediscache "zd/internal/cache/redis"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	PodId uint   `json:"pod_id"`
}

func main() {
	opts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	cache := rediscache.NewRedisCache(opts)
	defer cache.GracefulShutdown()

	// = = = Cache Example = = =
	exampleData := []User{
		{
			Id:    4,
			Name:  "Aldrick",
			PodId: 1,
		},
		{
			Id:    2,
			Name:  "Daniel",
			PodId: 3,
		},
		{
			Id:    2,
			Name:  "Daphne",
			PodId: 4,
		},
	}
	key := "userData"
	rawData, err := json.Marshal(exampleData)
	if err != nil {
		log.Fatalf("Failed to marshal data: %q", err)
	}
	ctx := context.Background()

	_, err = cache.CacheData(ctx, key, string(rawData), 1*time.Minute)
	if err != nil {
		log.Fatalf("Failed to cache data: %q", err)
	}

	result, err := cache.CheckCache(ctx, key)
	if err != nil {
		log.Fatalf("Cache miss: %q", err)
	}

	fmt.Printf("Data from cache: %+v\n", result)
}

func cachingJSON(client *redis.Client) {
	exampleData := []User{
		{
			Id:    4,
			Name:  "Aldrick",
			PodId: 1,
		},
		{
			Id:    2,
			Name:  "Daniel",
			PodId: 3,
		},
		{
			Id:    2,
			Name:  "Daphne",
			PodId: 4,
		},
	}

	ctx := context.Background()

	data, err := json.Marshal(exampleData)
	if err != nil {
		log.Fatalf("Error marshalling data: %q", err)
	}

	err = cacheSet(ctx, client, "jsonTest", data)
	if err != nil {
		log.Fatalf("Failed to cache data: %q", err)
	}

	cachedData, err := cacheGet(ctx, client, "jsonTest")
	if err != nil {
		log.Fatalf("Failed to get data: %q", err)
	}

	fmt.Printf("Got the data back: %+v\n", cachedData)

	time.Sleep(61 * time.Second)

	cachedData, err = cacheGet(ctx, client, "jsonTest")
	if err != nil {
		log.Fatalf("Failed to get data: %q", err)
	}

	fmt.Printf("Got the data back: %+v\n", cachedData)
}

func cacheSet(ctx context.Context, client *redis.Client, key string, data []byte) error {
	_, err := client.Set(ctx, key, data, 1*time.Minute).Result()
	return err
}

func cacheGet(ctx context.Context, client *redis.Client, key string) ([]User, error) {
	returnedData := client.Get(ctx, key)
	if returnedData == nil {
		return nil, fmt.Errorf("failed to get data")
	}

	dataBytes, err := returnedData.Bytes()
	if err != nil {
		log.Fatalf("Failed to convert to bytes: %q\n", dataBytes)
	}
	var umData []User
	err = json.Unmarshal(dataBytes, &umData)
	if err != nil {
		log.Fatalf("Failed to unmarshal the json data: %q", err)
	}

	return umData, nil
}
