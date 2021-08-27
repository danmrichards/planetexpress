package events

import (
	"context"
	"fmt"

	"github.com/danmrichards/planetexpress/internal/event"
	"github.com/go-redis/redis/v8"
)

// RedisService is a handler.EventService backed by Redis.
type RedisService struct {
	rc     *redis.Client
	stream string
}

// NewRedisService returns a new service instantiated with the given client
// and event stream name.
func NewRedisService(c *redis.Client, stream string) *RedisService {
	return &RedisService{rc: c, stream: stream}
}

// PackageAllocate implements handler.EventService.
func (r RedisService) PackageAllocate(ctx context.Context, id string, size int) error {
	evt := event.PackageEvent{
		ID:   id,
		Typ:  event.PackageAllocate,
		Size: size,
	}

	args := &redis.XAddArgs{
		Stream: r.stream,
		Values: evt.Values(),
	}
	if _, err := r.rc.XAdd(ctx, args).Result(); err != nil {
		return fmt.Errorf("could not dispatch event: %w", err)
	}

	return nil
}
