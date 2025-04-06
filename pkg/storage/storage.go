package storage

import (
	"context"
	"time"
)

type Storage interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	WithNamespace(namespace string) Storage
}
