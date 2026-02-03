package observability

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MetricsMiddleware records HTTP request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start)
		RecordHTTPRequest(
			c.Request.Context(),
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			duration,
		)
	}
}

// LoggingMiddleware logs HTTP requests and responses with structured logging
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID for correlation
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore body for handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response writer wrapper to capture response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Log request
		logger.InfoContext(c.Request.Context(),
			"HTTP Request",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("payload", string(requestBody)),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		logLevel := slog.LevelInfo
		if c.Writer.Status() >= 400 {
			logLevel = slog.LevelError
		}

		logger.Log(c.Request.Context(), logLevel,
			"HTTP Response",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status_code", c.Writer.Status()),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("response_body", blw.body.String()),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.ErrorContext(c.Request.Context(),
					"Request Error",
					slog.String("request_id", requestID),
					slog.String("error", err.Error()),
				)
			}
		}
	}
}

// bodyLogWriter wraps gin.ResponseWriter to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
