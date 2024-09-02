package repo

import (
	"context"
	"time"
)

// Redis 操作接口层

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
