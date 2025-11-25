package status

import (
	"actions-service/internal/models"
	"context"
)

type WorkcenterPort interface {
	GetWorkcenterDTO(ctx context.Context, id string) (*models.WorkcenterDTO, error)
}

type OperatorPort interface {
	ClockOut(ctx context.Context, operatorID, workcenterID string) error
}

