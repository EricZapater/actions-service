package workcenter

import (
	"actions-service/internal/models"
	"context"
)

type StatusPort interface {
	GetDefaultStatus(ctx context.Context) (models.StatusDTO, error)
	//StatusIn(ctx context.Context, workcenterID, statusID string, reasonID *string) error
}