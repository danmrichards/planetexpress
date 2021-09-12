package events

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danmrichards/planetexpress/internal/event"
	"github.com/go-redis/redis/v8"
)

// eventCallbackFunc is a callback function called when a new event is picked
// up by the listener.
type eventCallbackFunc func(*event.PackageEvent) error

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
func (r *RedisService) PackageAllocate(ctx context.Context, id string, size int) error {
	evt := event.PackageEvent{
		PackageID:   id,
		Typ:         event.PackageAllocate,
		PackageSize: size,
	}

	if _, err := r.rc.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: evt.Values(),
	}).Result(); err != nil {
		return fmt.Errorf("could not dispatch event: %w", err)
	}

	return nil
}

// Listen listens for events on the given stream from the last given event.
//
// The callback function, f, is called for each received event.
func (r *RedisService) Listen(ctx context.Context, f eventCallbackFunc, lastEvt *event.PackageEvent) error {
	// Indicates the stream to listen for events from and the starting point.
	// By default, we use the special 0 ID which just starts listening from the
	// first event.
	streams := []string{r.stream, "0"}

	// If we know the details of the most recent event, start listening from
	// there instead, in case we missed anything.
	if lastEvt != nil {
		streams[1] = lastEvt.EventID
	}

	for {
		res, err := r.rc.XRead(ctx, &redis.XReadArgs{
			Streams: streams,
			Count:   1,
		}).Result()

		switch {
		case errors.Is(err, context.Canceled):
			return nil
		case err != nil:
			return fmt.Errorf("listen for events: %w", err)
		}

		// Next iteration should start from the ID just seen.
		streams[1] = res[0].Messages[0].ID

		evt, err := event.FromValues(
			res[0].Messages[0].ID, res[0].Messages[0].Values,
		)
		if err != nil {
			log.Printf("parse event: %v", err)
			continue
		}

		if err = f(evt); err != nil {
			log.Printf("event callback: %v", err)
		}
	}
}
