package logging

import (
	"io"
	"log/slog"
	"os"
)

func New(verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	h := slog.NewTextHandler(io.Writer(os.Stdout), &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
