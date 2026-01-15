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
	
	pipe := r.client.Pipeline()
	pipe.Set(ctx, "status:"+id, data, 0)
	pipe.SAdd(ctx, "statuses", id)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.StatusDTO, error) {
	data, err := r.client.Get(ctx, "status:"+id).Bytes()
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
	// Using consistent workcenter namespace
	if err := r.client.Set(ctx, "workcenter:"+id, data, 0).Err(); err != nil {
		return err
	}
	return nil
}


