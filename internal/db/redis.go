package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, addr string, dbNum int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   dbNum,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}