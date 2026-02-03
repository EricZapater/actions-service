package observability

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	meterProvider  *metric.MeterProvider
	loggerProvider *sdklog.LoggerProvider
)

// Config holds the observability configuration
type Config struct {
	OtelEndpoint   string
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// InitTelemetry initializes OpenTelemetry with metrics and logs support
func InitTelemetry(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OtelEndpoint),
		otlpmetricgrpc.WithInsecure(), // Use insecure for local Alloy
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Create meter provider
	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(10*time.Second), // Export every 10 seconds
		)),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create OTLP log exporter
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.OtelEndpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	// Create logger provider
	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)

	log.Printf("OpenTelemetry initialized: service=%s, endpoint=%s (metrics + logs)", cfg.ServiceName, cfg.OtelEndpoint)

	// Return shutdown function
	return func(ctx context.Context) error {
		log.Println("Shutting down OpenTelemetry...")
		if err := meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
		if err := loggerProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown logger provider: %w", err)
		}
		return nil
	}, nil
}

// GetLoggerProvider returns the global logger provider
func GetLoggerProvider() *sdklog.LoggerProvider {
	return loggerProvider
}

// GetMeter returns the global meter for creating metrics
func GetMeter() {
	otel.Meter("actions-service")
}
