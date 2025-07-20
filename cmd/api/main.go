package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/bootstrap"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/system/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	ctx := context.Background()

	if os.Getenv("MIGRATE_ONLY") == "1" {
		app, err := bootstrap.Build(ctx, cfg)
		if err != nil {
			log.Fatalf("migrate-only build: %v", err)
		}
		app.DB.Close()
		log.Println("migrations applied (MIGRATE_ONLY=1), exiting")
		return
	}

	app, err := bootstrap.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := app.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	app.Shutdown(ctxShutdown)
	log.Println("shutdown complete")
}
