package observability

import (
	"context"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	// HTTP metrics
	httpRequestDuration metric.Float64Histogram
	httpRequestsTotal   metric.Int64Counter

	// Business metrics
	operatorClockInTotal    metric.Int64Counter
	operatorClockOutTotal   metric.Int64Counter
	statusChangesTotal      metric.Int64Counter
	workOrderPhaseInTotal   metric.Int64Counter
	workOrderPhaseOutTotal  metric.Int64Counter
	shiftChangesTotal       metric.Int64Counter
	shiftChangeDuration     metric.Float64Histogram

	metricsOnce sync.Once
)

// InitMetrics initializes all application metrics
func InitMetrics() error {
	var err error
	metricsOnce.Do(func() {
		meter := otel.Meter("actions-service")

		// HTTP metrics
		httpRequestDuration, err = meter.Float64Histogram(
			"http_request_duration_ms",
			metric.WithDescription("HTTP request duration in milliseconds"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			log.Printf("Failed to create http_request_duration_ms metric: %v", err)
			return
		}

		httpRequestsTotal, err = meter.Int64Counter(
			"http_requests_total",
			metric.WithDescription("Total number of HTTP requests"),
		)
		if err != nil {
			log.Printf("Failed to create http_requests_total metric: %v", err)
			return
		}

		// Business metrics
		operatorClockInTotal, err = meter.Int64Counter(
			"operator_clockin_total",
			metric.WithDescription("Total number of operator clock-ins"),
		)
		if err != nil {
			log.Printf("Failed to create operator_clockin_total metric: %v", err)
			return
		}

		operatorClockOutTotal, err = meter.Int64Counter(
			"operator_clockout_total",
			metric.WithDescription("Total number of operator clock-outs"),
		)
		if err != nil {
			log.Printf("Failed to create operator_clockout_total metric: %v", err)
			return
		}

		statusChangesTotal, err = meter.Int64Counter(
			"status_changes_total",
			metric.WithDescription("Total number of status changes"),
		)
		if err != nil {
			log.Printf("Failed to create status_changes_total metric: %v", err)
			return
		}

		workOrderPhaseInTotal, err = meter.Int64Counter(
			"workorderphase_in_total",
			metric.WithDescription("Total number of work order phase ins (load)"),
		)
		if err != nil {
			log.Printf("Failed to create workorderphase_in_total metric: %v", err)
			return
		}

		workOrderPhaseOutTotal, err = meter.Int64Counter(
			"workorderphase_out_total",
			metric.WithDescription("Total number of work order phase outs (unload)"),
		)
		if err != nil {
			log.Printf("Failed to create workorderphase_out_total metric: %v", err)
			return
		}

		shiftChangesTotal, err = meter.Int64Counter(
			"shift_changes_total",
			metric.WithDescription("Total number of shift changes"),
		)
		if err != nil {
			log.Printf("Failed to create shift_changes_total metric: %v", err)
			return
		}

		shiftChangeDuration, err = meter.Float64Histogram(
			"shift_change_duration_ms",
			metric.WithDescription("Duration of shift change operations in milliseconds"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			log.Printf("Failed to create shift_change_duration_ms metric: %v", err)
			return
		}

		log.Println("All metrics initialized successfully")
	})
	return err
}

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	if httpRequestDuration != nil {
		httpRequestDuration.Record(ctx, float64(duration.Milliseconds()),
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}

	if httpRequestsTotal != nil {
		httpRequestsTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}
}

// RecordOperatorClockIn records an operator clock-in event
func RecordOperatorClockIn(ctx context.Context, operatorID, workcenterID string) {
	if operatorClockInTotal != nil {
		operatorClockInTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("operator_id", operatorID),
				attribute.String("workcenter_id", workcenterID),
			),
		)
	}
}

// RecordOperatorClockOut records an operator clock-out event
func RecordOperatorClockOut(ctx context.Context, operatorID, workcenterID string) {
	if operatorClockOutTotal != nil {
		operatorClockOutTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("operator_id", operatorID),
				attribute.String("workcenter_id", workcenterID),
			),
		)
	}
}

// RecordStatusChange records a status change event
func RecordStatusChange(ctx context.Context, workcenterID, statusID string) {
	if statusChangesTotal != nil {
		statusChangesTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("workcenter_id", workcenterID),
				attribute.String("status_id", statusID),
			),
		)
	}
}

// RecordWorkOrderPhaseIn records a work order phase in (load) event
func RecordWorkOrderPhaseIn(ctx context.Context, workcenterID string) {
	if workOrderPhaseInTotal != nil {
		workOrderPhaseInTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("workcenter_id", workcenterID),
			),
		)
	}
}

// RecordWorkOrderPhaseOut records a work order phase out (unload) event
func RecordWorkOrderPhaseOut(ctx context.Context, workcenterID string) {
	if workOrderPhaseOutTotal != nil {
		workOrderPhaseOutTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("workcenter_id", workcenterID),
			),
		)
	}
}

// RecordShiftChange records a shift change event with duration
func RecordShiftChange(ctx context.Context, workcenterID, shiftDetailID string, duration time.Duration) {
	if shiftChangesTotal != nil {
		shiftChangesTotal.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("workcenter_id", workcenterID),
				attribute.String("shift_detail_id", shiftDetailID),
			),
		)
	}

	if shiftChangeDuration != nil {
		shiftChangeDuration.Record(ctx, float64(duration.Milliseconds()),
			metric.WithAttributes(
				attribute.String("workcenter_id", workcenterID),
				attribute.String("shift_detail_id", shiftDetailID),
			),
		)
	}
}
