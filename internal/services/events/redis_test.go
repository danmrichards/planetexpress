package events

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/danmrichards/planetexpress/internal/event"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
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
		"package_id": "foo",
		"type":       string(event.PackageAllocate),
		"size":       "10",
	}, res[0].Values)
}
