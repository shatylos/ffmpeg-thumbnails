package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/shatylos/ffmpeg-thumbnails/tools/apperrors"
	"github.com/shatylos/ffmpeg-thumbnails/tools/logger"
)

// StartServer runs an HTTP server that serves the latest screenshots from
// storage. It shuts down when ctx is cancelled. It is a no-op when the
// server is disabled in config.
func StartServer(ctx context.Context, config Config, storage Storage) {
	if !config.ServerEnabled {
		return
	}

	srv := &http.Server{
		Addr:    config.ServerAddr,
		Handler: screenshotHandler(storage),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.PrintError(apperrors.Wrap(err, "http server shutdown"))
		}
	}()

	go func() {
		logger.Info("http server listening on " + config.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.PrintError(apperrors.Wrap(err, "http server failed"))
		}
	}()
}

// screenshotHandler serves the latest screenshot for the requested stream.
// The stream is identified by the request path, e.g. GET /cam1.jpg.
func screenshotHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := strings.TrimPrefix(r.URL.Path, "/")
		data, ok := storage.Get(output)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(data)
	}
}
