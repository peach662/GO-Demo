package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func OpenRedis(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}
