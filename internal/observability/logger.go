package observability

import (
	"context"
	"log/slog"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewLogger creates a new structured logger that exports via OTLP
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

	// Get the logger provider (must be initialized first)
	provider := GetLoggerProvider()
	if provider == nil {
		// Fallback to stdout if OTLP not initialized
		return slog.Default()
	}

	// Create OTLP handler for structured logging
	handler := otelslog.NewHandler("actions-service", otelslog.WithLoggerProvider(provider))
	
	// Wrap with level filter
	levelHandler := &levelFilterHandler{
		handler: handler,
		level:   level,
	}
	
	logger := slog.New(levelHandler)

	return logger
}

// levelFilterHandler wraps a handler to filter by log level
type levelFilterHandler struct {
	handler slog.Handler
	level   slog.Level
}

func (h *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *levelFilterHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{
		handler: h.handler.WithAttrs(attrs),
		level:   h.level,
	}
}

func (h *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{
		handler: h.handler.WithGroup(name),
		level:   h.level,
	}
}
