package locker

import (
"context"
"github.com/danakum/go-util/logger"
"github.com/danakum/go-util/redis"
)

var lockPath = `locks:`

func AquireLock(ctx context.Context, resourceType string, identifire string, value string) error {
	err := redis.Set(ctx, lockPath+resourceType+redis.KeySeparator+value, true)
	if err != nil {
		logger.Log().ErrorContext(ctx, err)
	}
	return err
}

func ReleaseLock(ctx context.Context, resourceType string, identifire string, value string) error {
	err := redis.Set(ctx, lockPath+resourceType+redis.KeySeparator+value, false)
	if err != nil {
		logger.Log().ErrorContext(ctx, err)
	}
	return err
}
