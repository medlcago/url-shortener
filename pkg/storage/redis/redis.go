package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
	"url-shortener/pkg/storage"
)

type Storage struct {
	client    *redis.Client
	namespace string
}

func New(cfg ...Config) storage.Storage {
	c := defaultConfig(cfg...)
	return &Storage{client: c.Client}

}

func (r *Storage) WithNamespace(namespace string) storage.Storage {
	if r.namespace != "" {
		namespace = r.namespace + "_" + namespace
	}
	return &Storage{
		client:    r.client,
		namespace: namespace,
	}
}

func (r *Storage) makeKey(key string) string {
	if r.namespace != "" {
		return r.namespace + ":" + key
	}
	return key
}

func (r *Storage) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Storage) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, r.makeKey(key)).Result()
}

func (r *Storage) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, r.makeKey(key)).Result()
	return exists > 0, err
}

func (r *Storage) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.makeKey(key)).Err()
}

func (r *Storage) Close() error {
	return r.client.Close()
}
