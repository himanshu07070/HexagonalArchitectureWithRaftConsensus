package database

import (
	"context"
	"fmt"
	"hexagonal-architecture/internal"

	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter() internal.DatabasePort {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisAdapter{client: rdb}
}

// Store stores the fileName-fileSize mapping in Redis.
func (r *RedisAdapter) Store(fileName string, fileSize int64) error {
	key := fmt.Sprintf("file:%s", fileName)
	return r.client.Set(context.Background(), key, fileSize, 0).Err()
}

// Retrieve retrieves the fileSize for a given fileName from Redis.
func (r *RedisAdapter) Retrieve(fileName string) (int64, error) {
	key := fmt.Sprintf("file:%s", fileName)
	val, err := r.client.Get(context.Background(), key).Int64()
	if err == redis.Nil {
		return 0, fmt.Errorf("File not found in Redis")
	}
	return val, err
}
