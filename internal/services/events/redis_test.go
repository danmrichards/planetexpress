package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/danmrichards/planetexpress/internal/event"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisService_PackageAllocate(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rc.Close()

	stream := "test"
	svc := NewRedisService(rc, stream)

	assert.NoError(t, svc.PackageAllocate(context.Background(), "foo", 10))

	res, err := rc.XRange(context.Background(), stream, "-", "+").Result()
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, map[string]interface{}{
		"package_id":   "foo",
		"type":         string(event.PackageAllocate),
		"package_size": "10",
	}, res[0].Values)
}

func TestRedisService_Listen(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rc.Close()

	stream := "test"
	svc := NewRedisService(rc, stream)

	assert.NoError(t, svc.PackageAllocate(context.Background(), "foo", 10))

	var called bool
	testCallback := func(evt *event.PackageEvent) error {
		if evt.Typ != event.PackageAllocate {
			return fmt.Errorf("unexpected type: %v", evt.Typ)
		}
		if evt.PackageID != "foo" {
			return fmt.Errorf("unexpected package ID: %q", evt.PackageID)
		}
		if evt.PackageSize != 10 {
			return fmt.Errorf("unexpected package size: %d", evt.PackageSize)
		}

		called = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, cancel)

	require.NoError(t, svc.Listen(ctx, testCallback, nil))
	require.True(t, called)
}
