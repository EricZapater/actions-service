package status

import (
	"actions-service/internal/models"
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisRepo *RedisRepo
}

func NewStatusRepository(client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	return &Repository{
		redisRepo: redisRepo,
	}
}

func (r *Repository) Set(ctx context.Context, id string, value models.StatusDTO) error {
	return r.redisRepo.Set(ctx, id, value)
}

func (r *Repository) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
    return r.redisRepo.SetWorkcenterDTO(ctx, id, value)
}

func (r *Repository) FindByID(ctx context.Context, id string) (models.StatusDTO, models.DataSource, error) {
	status, err := r.redisRepo.FindByID(ctx, id)
	if err == nil {
		return status, models.SourceRedis, nil
	}
	return models.StatusDTO{}, models.SourceNone, err
}


