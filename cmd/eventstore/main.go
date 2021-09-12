package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/danmrichards/planetexpress/internal/services/events"
	"github.com/danmrichards/planetexpress/internal/services/package_events"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	redisAddr   string
	pgqlAddr    string
	pgsqlDB     string
	pgsqlUser   string
	pgsqlPass   string
	eventStream string
)

func main() {
	flag.StringVar(&redisAddr, "redis-addr", "127.0.0.1:6379", "the ip:port of the Redis server")
	flag.StringVar(&pgqlAddr, "pgsql-addr", "127.0.0.1:5432", "the ip:port of the PostgreSQL server")
	flag.StringVar(&pgsqlDB, "pgsql-db", "", "the name of the PostgreSQL database")
	flag.StringVar(&pgsqlUser, "pgsql-user", "", "the PostgreSQL user")
	flag.StringVar(&pgsqlPass, "pgsql-pass", "", "the PostgreSQL password")
	flag.StringVar(&eventStream, "event-stream", "packages", "the name of the event stream")
	flag.Parse()

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

	pkgEvtSvc := package_events.NewSQLService(db)

	// We need the latest package event to use as the starting point for the
	// stream listener.
	pkgEvt, err := pkgEvtSvc.Latest()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}

	rc := redis.NewClient(&redis.Options{Addr: redisAddr})
	evtSvc := events.NewRedisService(rc, eventStream)

	ctx := signals.SetupSignalHandler()

	log.Printf("listening for events...")
	if err = evtSvc.Listen(ctx, pkgEvtSvc.Save, pkgEvt); err != nil {
		log.Fatalf("could not start API server: %v", err)
	}

	rc.Close() //nolint:errcheck
	db.Close() //nolint:errcheck
}
