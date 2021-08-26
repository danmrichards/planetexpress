package main

import (
	"flag"
	"log"
	"time"

	"github.com/danmrichards/planetexpress/internal/api/handler"
	"github.com/danmrichards/planetexpress/internal/http"
	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var bind string

const shutdownTimeout = 5 * time.Second

func main() {
	flag.StringVar(&bind, "", "127.0.0.1:5000", "the ip:port to bind the API server to")
	flag.Parse()

	r := mux.NewRouter()
	if err := handler.Init(r); err != nil {
		log.Fatalf("could not setup API handler: %v", err)
	}

	srv := http.NewServer(bind, shutdownTimeout, r)

	ctx := signals.SetupSignalHandler()

	log.Printf("starting api server on %q\n", bind)
	if err := srv.Serve(ctx); err != nil {
		log.Fatalf("could not start API server: %v", err)
	}
}
