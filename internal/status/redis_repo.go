package status

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

func (r *RedisRepo) Set(ctx context.Context, id string, value models.StatusDTO) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, id, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.StatusDTO, error) {
	data, err := r.client.Get(ctx, id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.StatusDTO{}, ErrStatusNotFound
		}
		return models.StatusDTO{}, err
	}
	var status models.StatusDTO
	if err := json.Unmarshal(data, &status); err != nil {
		return models.StatusDTO{}, err
	}
	return status, nil
}

func (r *RedisRepo) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    if err := r.client.Set(ctx, id, data, 0).Err(); err != nil {
        return err
    }
    return nil
}


