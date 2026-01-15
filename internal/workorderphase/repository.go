package workorderphase

import (
	"actions-service/internal/models"
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisRepo *RedisRepo
}

func NewWorkOrderPhaseRepository(client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	return &Repository{
		redisRepo: redisRepo,
	}
}

func (r *Repository) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
	return r.redisRepo.SetWorkcenterDTO(ctx, id, value)
}

func (r *Repository) List(ctx context.Context) ([]models.WorkcenterDTO, error) {
	return r.redisRepo.List(ctx)
}
