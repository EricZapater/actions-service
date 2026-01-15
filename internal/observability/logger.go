package observability

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a new structured logger with JSON output
func NewLogger(logLevel string) *slog.Logger {
	// Parse log level
	level := slog.LevelInfo
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	// Create JSON handler for structured logging
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
