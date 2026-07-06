package main

import (
	"context"
	"github.com/shatylos/ffmpeg-screenshots/internal/app"
	"github.com/shatylos/ffmpeg-screenshots/tools/logger"
	"os/signal"
	"syscall"
)

func main() {
	config, err := app.GetConfig()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.HandleQueue(ctx, config)

	<-ctx.Done()
	logger.Info("shutting down...")
}
