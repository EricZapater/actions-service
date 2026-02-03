package validator

import (
	"actions-service/internal/models"
	"context"
)

// ⭐ PORTS - Interfaces that validator uses to interact with other services
// These are defined HERE in the validator package to avoid circular dependencies

// OperatorPort defines operations the validator needs from the operator service
type OperatorPort interface {
	ClockOut(ctx context.Context, operatorID, workcenterID string) error
}

// StatusRepository defines operations the validator needs to get status information
type StatusRepository interface {
	FindByID(ctx context.Context, key string) (models.StatusDTO, models.DataSource, error)
}

// WorkOrderPhasePort defines operations the validator needs from the workorderphase service
type WorkOrderPhasePort interface {
	ForcePhaseOut(ctx context.Context, workcenterID, workorderPhaseID string) error
}

// WorkcenterPort defines operations the validator needs from the workcenter service
type WorkcenterPort interface {
	GetWorkcenterDTO(ctx context.Context, workcenterID string) (*models.WorkcenterDTO, error)
}
