package observability

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	meterProvider *metric.MeterProvider
)

// Config holds the observability configuration
type Config struct {
	OtelEndpoint   string
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// InitTelemetry initializes OpenTelemetry with metrics support
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

	log.Printf("OpenTelemetry initialized: service=%s, endpoint=%s", cfg.ServiceName, cfg.OtelEndpoint)

	// Return shutdown function
	return func(ctx context.Context) error {
		log.Println("Shutting down OpenTelemetry...")
		if err := meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
		return nil
	}, nil
}

// GetMeter returns the global meter for creating metrics
func GetMeter() {
	otel.Meter("actions-service")
}
