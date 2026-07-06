package main

import (
	"context"
	"github.com/shatylos/ffmpeg-thumbnails/internal/app"
	"github.com/shatylos/ffmpeg-thumbnails/tools/logger"
	"os/signal"
	"syscall"
)

func main() {
	config, err := app.GetConfig()
	if err != nil {
		panic(err)
	}

	storage, err := app.NewStorage(config)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.HandleQueue(ctx, config, storage)
	app.StartServer(ctx, config, storage)

	<-ctx.Done()
	logger.Info("shutting down...")
}
