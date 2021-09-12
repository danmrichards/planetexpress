package main

import (
	"flag"
	"log"
	"time"

	"github.com/danmrichards/planetexpress/internal/api/handler"
	"github.com/danmrichards/planetexpress/internal/http"
	"github.com/danmrichards/planetexpress/internal/services/events"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	bind        string
	redisAddr   string
	eventStream string
)

const shutdownTimeout = 5 * time.Second

func main() {
	flag.StringVar(&bind, "bind", "127.0.0.1:5000", "the ip:port to bind the API server to")
	flag.StringVar(&redisAddr, "redis-addr", "127.0.0.1:6379", "the ip:port of the Redis server")
	flag.StringVar(&eventStream, "event-stream", "packages", "the name of the event stream")
	flag.Parse()

	r := mux.NewRouter()
	rc := redis.NewClient(&redis.Options{Addr: redisAddr})
	evtSvc := events.NewRedisService(rc, eventStream)

	if err := handler.Init(r, evtSvc); err != nil {
		log.Fatalf("could not setup API handler: %v", err)
	}

	srv := http.NewServer(bind, shutdownTimeout, r)

	ctx := signals.SetupSignalHandler()

	log.Printf("starting api server on %q\n", bind)
	if err := srv.Serve(ctx); err != nil {
		log.Fatalf("could not start API server: %v", err)
	}

	rc.Close() //nolint:errcheck
}
