package workcenter

import (
	"actions-service/internal/models"
	"context"

	"github.com/redis/go-redis/v9"
)



type Repository struct {
	redisRepo *RedisRepo
}

func NewWorkcenterRepository(client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	return &Repository{
		redisRepo: redisRepo,
	}
}

func(r *Repository) Set(ctx context.Context, id string, value models.WorkcenterDTO)error{
	return r.redisRepo.Set(ctx, id, value)
}

func (r *Repository) FindByID(ctx context.Context, id string) (models.WorkcenterDTO, models.DataSource, error) {
	workcenter, err := r.redisRepo.FindByID(ctx, id)
	if err == nil {
		return workcenter, models.SourceRedis, nil
	}
	return models.WorkcenterDTO{}, models.SourceNone, err
}

func(r *Repository) List(ctx context.Context) ([]models.WorkcenterDTO, error){
	return r.redisRepo.List(ctx)
}

func (r *Repository) Delete(ctx context.Context, id string) error{
	return r.redisRepo.Delete(ctx, id)
}