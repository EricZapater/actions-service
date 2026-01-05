package operator

import (
	"actions-service/internal/models"
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisRepo *RedisRepo
}

func NewOperatorRepository(client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	return &Repository{
		redisRepo: redisRepo,
	}
}

func(r *Repository) Set(ctx context.Context, id string, value models.OperatorDTO)error{
	return r.redisRepo.Set(ctx, id, value)
}

func(r *Repository) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO)error{
	return r.redisRepo.SetWorkcenterDTO(ctx, id, value)
}

func (r *Repository) FindByID(ctx context.Context, id string) (models.OperatorDTO, models.DataSource, error){
	operator, err := r.redisRepo.FindByID(ctx, id)
	if err == nil {
		return operator, models.SourceRedis, nil
	}
	return models.OperatorDTO{}, models.SourceNone, err
}