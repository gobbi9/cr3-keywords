package pipeline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func extractPreviewImage(ctx context.Context, cr3Path string, forceExif bool) ([]byte, error) {
	if forceExif {
		return extractPreviewImageWithExifTool(ctx, cr3Path)
	}

	preview, pureErr := extractPreviewImagePureGo(cr3Path)
	if pureErr == nil {
		return preview, nil
	}

	preview, exiftoolErr := extractPreviewImageWithExifTool(ctx, cr3Path)
	if exiftoolErr == nil {
		return preview, nil
	}

	return nil, fmt.Errorf("extract preview failed (pure-go: %v, exiftool fallback: %w)", pureErr, exiftoolErr)
}

func extractPreviewImagePureGo(cr3Path string) ([]byte, error) {
	data, err := os.ReadFile(cr3Path)
	if err != nil {
		return nil, fmt.Errorf("read CR3 file: %w", err)
	}

	const (
		jpegSOI0 = 0xFF
		jpegSOI1 = 0xD8
		jpegEOI0 = 0xFF
		jpegEOI1 = 0xD9
		minJPEG  = 64 * 1024
	)

	var best []byte
	searchFrom := 0
	for searchFrom < len(data)-3 {
		startRel := bytes.Index(data[searchFrom:], []byte{jpegSOI0, jpegSOI1, 0xFF})
		if startRel < 0 {
			break
		}
		start := searchFrom + startRel

		endRel := bytes.Index(data[start+2:], []byte{jpegEOI0, jpegEOI1})
		if endRel < 0 {
			break
		}
		end := start + 2 + endRel + 2

		candidate := data[start:end]
		if len(candidate) >= minJPEG && len(candidate) > len(best) {
			best = candidate
		}

		searchFrom = start + 3
	}

	if len(best) == 0 {
		return nil, errors.New("no embedded JPEG preview found")
	}

	out := make([]byte, len(best))
	copy(out, best)
	return out, nil
}

func extractPreviewImageWithExifTool(ctx context.Context, cr3Path string) ([]byte, error) {
	if _, err := exec.LookPath("exiftool"); err != nil {
		return nil, fmt.Errorf("exiftool not found in PATH")
	}

	cmd := exec.CommandContext(ctx, "exiftool", "-b", "-PreviewImage", cr3Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	b, readErr := io.ReadAll(stdout)
	waitErr := cmd.Wait()
	if readErr != nil {
		return nil, readErr
	}
	if waitErr != nil {
		return nil, waitErr
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("empty preview image stream")
	}

	return b, nil
}

func readOrientationFromCR3(ctx context.Context, cr3Path string) (int, error) {
	if _, err := exec.LookPath("exiftool"); err != nil {
		return 0, fmt.Errorf("exiftool not found in PATH")
	}

	cmd := exec.CommandContext(ctx, "exiftool", "-Orientation", "-n", "-s", "-s", "-s", cr3Path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	s := bytes.TrimSpace(out)
	if len(s) == 0 {
		return 0, errors.New("empty orientation output")
	}

	var orientation int
	if _, err := fmt.Sscanf(string(s), "%d", &orientation); err != nil {
		return 0, fmt.Errorf("parse orientation: %w", err)
	}
	if orientation < 1 || orientation > 8 {
		return 0, fmt.Errorf("invalid orientation value: %d", orientation)
	}

	return orientation, nil
}
