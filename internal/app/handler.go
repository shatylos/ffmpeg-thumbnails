package app

import (
	"context"
	"fmt"
	"github.com/shatylos/ffmpeg-screenshots/tools/apperrors"
	"github.com/shatylos/ffmpeg-screenshots/tools/logger"
	"os/exec"
	"path/filepath"
	"time"
)

var inProcess = map[string]bool{}

func HandleQueue(ctx context.Context, config Config) {
	queue := make(chan StreamConfig, config.Forks)

	go func() {
		for stream := range queue {
			go handleStream(ctx, config, stream)
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

func handleStream(ctx context.Context, config Config, streamConfig StreamConfig) {
	if inProcess[streamConfig.Output] {
		logger.Warning(fmt.Sprintf("stream handling [%s] still in progress", streamConfig.Output))
		return
	}
	inProcess[streamConfig.Output] = true
	defer func() {
		inProcess[streamConfig.Output] = false
	}()

	output := filepath.Join(config.Outputdir, streamConfig.Output)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", streamConfig.Src,
		"-frames:v", "1",
		"-update", "true",
		"-y", output,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		logger.PrintError(apperrors.Wrap(err, "ffmpeg failed for %s: %s", streamConfig.Src, string(out)))
		return
	}
}
