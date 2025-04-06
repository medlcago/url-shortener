package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
	"url-shortener/config"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.RedisAddr,
		Password:     cfg.Redis.RedisPassword,
		DB:           cfg.Redis.RedisDb,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Second * cfg.Redis.PoolTimeout,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
