package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/danmrichards/planetexpress/internal/api/handler"
	"github.com/danmrichards/planetexpress/internal/http"
	"github.com/danmrichards/planetexpress/internal/services/events"
	"github.com/danmrichards/planetexpress/internal/services/ship"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	bind        string
	redisAddr   string
	pgqlAddr    string
	pgsqlDB     string
	pgsqlUser   string
	pgsqlPass   string
	eventStream string
)

const shutdownTimeout = 5 * time.Second

func main() {
	flag.StringVar(&bind, "bind", "127.0.0.1:5000", "the ip:port to bind the API server to")
	flag.StringVar(&redisAddr, "redis-addr", "127.0.0.1:6379", "the ip:port of the Redis server")
	flag.StringVar(&pgqlAddr, "pgsql-addr", "127.0.0.1:5432", "the ip:port of the PostgreSQL server")
	flag.StringVar(&pgsqlDB, "pgsql-db", "", "the name of the PostgreSQL database")
	flag.StringVar(&pgsqlUser, "pgsql-user", "", "the PostgreSQL user")
	flag.StringVar(&pgsqlPass, "pgsql-pass", "", "the PostgreSQL password")
	flag.StringVar(&eventStream, "event-stream", "packages", "the name of the event stream")
	flag.Parse()

	r := mux.NewRouter()
	rc := redis.NewClient(&redis.Options{Addr: redisAddr})
	evtSvc := events.NewRedisService(rc, eventStream)

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			pgsqlUser, pgsqlPass, pgqlAddr, pgsqlDB,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shipSvc := ship.NewSQLService(db)

	if err = handler.Init(r, evtSvc, shipSvc); err != nil {
		log.Fatalf("could not setup API handler: %v", err)
	}

	srv := http.NewServer(bind, shutdownTimeout, r)

	ctx := signals.SetupSignalHandler()

	log.Printf("starting api server on %q\n", bind)
	if err = srv.Serve(ctx); err != nil {
		log.Fatalf("could not start API server: %v", err)
	}

	rc.Close() //nolint:errcheck
	db.Close() //nolint:errcheck
}
