package workorderphase

import (
	"actions-service/internal/models"
	"context"
)

type WorkcenterPort interface {
	GetWorkcenterDTO(ctx context.Context, id string) (*models.WorkcenterDTO, error)
}

type StatusPort interface {
	FindByID(ctx context.Context, workcenterID, statusID string) (models.StatusDTO, error)
}

type OperatorPort interface {
	ClockOut(ctx context.Context, operatorID, workcenterID string) error
}

