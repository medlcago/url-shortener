package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage struct {
	client    *redis.Client
	namespace string
}

func WithRedisClient(client *redis.Client) Storage {
	return &RedisStorage{client: client}
}
func (r *RedisStorage) WithNamespace(namespace string) Storage {
	if r.namespace != "" {
		namespace = r.namespace + "_" + namespace
	}
	return &RedisStorage{
		client:    r.client,
		namespace: namespace,
	}
}

func (r *RedisStorage) makeKey(key string) string {
	if r.namespace != "" {
		return r.namespace + ":" + key
	}
	return key
}

func (r *RedisStorage) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisStorage) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, r.makeKey(key)).Result()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, r.makeKey(key)).Result()
	return exists > 0, err
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.makeKey(key)).Err()
}
