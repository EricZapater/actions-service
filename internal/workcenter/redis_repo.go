package workcenter

import (
	"actions-service/internal/models"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		client: client,
	}
}

func(r *RedisRepo) Set(ctx context.Context, id string, value models.WorkcenterDTO) error{
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, id, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.WorkcenterDTO, error){
	data, err := r.client.Get(ctx, id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.WorkcenterDTO{}, ErrWorkcenterNotFound
		}
		return models.WorkcenterDTO{}, err
	}

	var workcenter models.WorkcenterDTO
	if err := json.Unmarshal(data, &workcenter); err != nil {
		return models.WorkcenterDTO{}, err
	}
	return workcenter, nil
}

func(r *RedisRepo) List(ctx context.Context) ([]models.WorkcenterDTO, error){
	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	var workcenters []models.WorkcenterDTO
	for _, key := range keys {
		workcenter, err := r.FindByID(ctx, key)
		if err != nil {
			return nil, err
		}
		workcenters = append(workcenters, workcenter)
	}
	return workcenters, nil
}