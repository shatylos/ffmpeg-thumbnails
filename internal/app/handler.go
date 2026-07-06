package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/shatylos/ffmpeg-thumbnails/tools/apperrors"
	"github.com/shatylos/ffmpeg-thumbnails/tools/logger"
	"os/exec"
	"sync"
	"time"
)

var (
	inProcessMu sync.Mutex
	inProcess   = map[string]bool{}
)

func HandleQueue(ctx context.Context, config Config, storage Storage) {
	queue := make(chan StreamConfig, config.Forks)

	go func() {
		for stream := range queue {
			go handleStream(ctx, config, stream, storage)
		}
	}()

	go func() {
		ticker := time.NewTicker(config.Frequency)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, streamConfig := range config.Streams {
					queue <- streamConfig
				}
			}
		}
	}()
}

func handleStream(ctx context.Context, config Config, streamConfig StreamConfig, storage Storage) {
	inProcessMu.Lock()
	if inProcess[streamConfig.Output] {
		inProcessMu.Unlock()
		logger.Warning(fmt.Sprintf("stream handling [%s] still in progress", streamConfig.Output))
		return
	}
	inProcess[streamConfig.Output] = true
	inProcessMu.Unlock()
	defer func() {
		inProcessMu.Lock()
		inProcess[streamConfig.Output] = false
		inProcessMu.Unlock()
	}()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", streamConfig.Src,
		"-frames:v", "1",
		"-f", "image2pipe",
		"-c:v", "mjpeg",
		"-",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.PrintError(apperrors.Wrap(err, "ffmpeg failed for %s: %s", streamConfig.Src, stderr.String()))
		return
	}

	if err := storage.Save(streamConfig.Output, stdout.Bytes()); err != nil {
		logger.PrintError(err)
		return
	}
}
