package core_dependencies

import (
	"context"
	"time"
)

type Cache interface {
	CacheData(ctx context.Context, key string, data string, duration time.Duration) (string, error)
	CheckCache(ctx context.Context, key string) (string, error)
	GracefulShutdown()
}
