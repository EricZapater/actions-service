package operator

import (
	"actions-service/internal/models"
	"context"
)

type WorkcenterPort interface {
	GetWorkcenterDTO(ctx context.Context, id string) (*models.WorkcenterDTO, error)
}