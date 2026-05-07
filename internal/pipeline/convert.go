package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func cr3ToJPGs(ctx context.Context, logger *slog.Logger, progress *ProgressBar, jpgDir string, files []string, dryRun bool, useExif bool) error {
	if len(files) == 0 {
		return fmt.Errorf("no CR3 files to convert")
	}

	if !dryRun {
		if err := os.MkdirAll(jpgDir, 0o755); err != nil {
			return fmt.Errorf("create jpg dir: %w", err)
		}
	}

	total := len(files)
	for idx, f := range files {
		base := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		out := filepath.Join(jpgDir, base+".jpg")

		if _, err := os.Stat(out); err == nil {
			logger.Debug("skipping CR3 because JPG exists", "file", f)
			progress.Render(idx+1, total)
			continue
		}

		if dryRun {
			logger.Debug("[dry-run] would generate JPG", "from", f, "to", out)
			progress.Render(idx+1, total)
			continue
		}

		preview, err := extractPreviewImage(ctx, f, useExif)
		if err != nil {
			logger.Error("failed to extract preview image", "file", f, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		img, err := imaging.Decode(bytes.NewReader(preview))
		if err != nil {
			logger.Error("failed to decode preview JPEG", "file", f, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		orientation := readOrientation(ctx, f, preview)
		img = applyOrientation(img, orientation)
		img = imaging.Fit(img, 2000, 2000, imaging.Lanczos)

		if err := imaging.Save(img, out, imaging.JPEGQuality(85)); err != nil {
			logger.Error("failed to write JPG", "file", out, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		progress.Render(idx+1, total)
	}

	return nil
}

func readOrientation(ctx context.Context, cr3Path string, preview []byte) int {
	if o, err := readOrientationFromCR3(ctx, cr3Path); err == nil {
		return o
	}

	x, err := exif.Decode(bytes.NewReader(preview))
	if err != nil {
		return 1
	}
	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 1
	}
	o, err := tag.Int(0)
	if err != nil {
		return 1
	}
	return o
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		// EXIF 6 means display image rotated 90° clockwise.
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		// EXIF 8 means display image rotated 270° clockwise (90° counter-clockwise).
		return imaging.Rotate90(img)
	default:
		return img
	}
}
