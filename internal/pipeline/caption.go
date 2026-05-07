package pipeline

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"cr3-keywords/internal/lm"
)

func batchCaption(ctx context.Context, logger *slog.Logger, progress *ProgressBar, client *lm.Client, jpgDir string, outputDir string, model string, promptFile string, images []string, dryRun bool) error {
	if len(images) == 0 {
		return fmt.Errorf("no JPG files to caption")
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
	for idx, img := range images {
		if !strings.HasPrefix(img, jpgDir) {
			logger.Debug("caption input image outside jpg dir", "image", img)
		}

		base := strings.TrimSuffix(filepath.Base(img), filepath.Ext(img))
		out := filepath.Join(outputDir, base+".txt")

		if _, err := os.Stat(out); err == nil {
			logger.Debug("skipping JPG because TXT exists", "file", img)
			progress.Render(idx+1, total)
			continue
		}

		if dryRun {
			logger.Debug("[dry-run] would caption image", "from", img, "to", out)
			progress.Render(idx+1, total)
			continue
		}

		imgBytes, err := os.ReadFile(img)
		if err != nil {
			logger.Error("failed to read JPG", "file", img, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		b64 := base64.StdEncoding.EncodeToString(imgBytes)
		caption, err := client.ChatCaption(ctx, model, prompt, b64)
		if err != nil {
			logger.Error("failed to caption image", "file", img, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		if err := os.WriteFile(out, []byte(strings.TrimSpace(caption)+"\n"), 0o644); err != nil {
			logger.Error("failed to write TXT", "file", out, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		progress.Render(idx+1, total)
	}

	return nil
}
