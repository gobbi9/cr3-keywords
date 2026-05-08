package logging

import (
	"io"
	"log/slog"
	"os"
)

// New returns a slog logger configured for stdout.
//
// When verbose is true, the logger uses debug level; otherwise it uses info level.
func New(verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	h := slog.NewTextHandler(io.Writer(os.Stdout), &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
