package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhohov/contest-service/judge-service/config"
	"github.com/zhohov/contest-service/judge-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
