package pipeline

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"cr3-keywords/internal/lm"
)

func batchKeywords(ctx context.Context, logger *slog.Logger, progress *ProgressBar, lmClient *lm.Client, jpgDir string, outputDir string, model string, promptFile string, images []string, geoByBase map[string]GeoMetadata, dryRun bool) error {
	if len(images) == 0 {
		return fmt.Errorf("no JPG files to keyword")
	}

	promptBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return fmt.Errorf("read prompt file: %w", err)
	}
	prompt := string(promptBytes)

	if !dryRun {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	total := len(images)
	started := time.Now()

	var renderedCurrent atomic.Int64
	var renderMu sync.Mutex
	render := func(current int) {
		suffix := fmt.Sprintf(" (%s)", time.Since(started).Truncate(time.Second))
		renderMu.Lock()
		defer renderMu.Unlock()
		progress.RenderWithSuffix(current, total, suffix)
	}

	render(0)
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				render(int(renderedCurrent.Load()))
			case <-done:
				return
			}
		}
	}()
	defer func() {
		ticker.Stop()
		close(done)
	}()

	for idx, img := range images {
		if !strings.HasPrefix(img, jpgDir) {
			logger.Debug("keyword input image outside jpg dir", "image", img)
		}

		base := strings.TrimSuffix(filepath.Base(img), filepath.Ext(img))
		out := filepath.Join(outputDir, base+".txt")

		if _, err := os.Stat(out); err == nil {
			logger.Debug("skipping JPG because TXT exists", "file", img)
			renderedCurrent.Store(int64(idx + 1))
			render(idx + 1)
			continue
		}

		if dryRun {
			logger.Debug("[dry-run] would caption image", "from", img, "to", out)
			renderedCurrent.Store(int64(idx + 1))
			render(idx + 1)
			continue
		}

		imgBytes, err := os.ReadFile(img)
		if err != nil {
			logger.Error("failed to read JPG", "file", img, "error", err)
			renderedCurrent.Store(int64(idx + 1))
			render(idx + 1)
			continue
		}

		effectivePrompt := prompt
		if geo, ok := geoByBase[base]; ok {
			effectivePrompt = prependGeoContext(prompt, geo)
		}
		logger.Debug("sending prompt to LM Studio", "server", "http://localhost:1234", "image", img, "prompt", effectivePrompt)

		b64 := base64.StdEncoding.EncodeToString(imgBytes)
		keywordsCaption, err := lmClient.PromptWithImage(ctx, model, effectivePrompt, b64)
		if err != nil {
			logger.Error("failed to keyword image", "file", img, "error", err)
			renderedCurrent.Store(int64(idx + 1))
			render(idx + 1)
			continue
		}

		if err := os.WriteFile(out, []byte(strings.TrimSpace(keywordsCaption)+"\n"), 0o644); err != nil {
			logger.Error("failed to write TXT", "file", out, "error", err)
			renderedCurrent.Store(int64(idx + 1))
			render(idx + 1)
			continue
		}

		renderedCurrent.Store(int64(idx + 1))
		render(idx + 1)
	}

	return nil
}
