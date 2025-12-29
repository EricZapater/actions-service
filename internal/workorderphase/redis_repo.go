package workorderphase

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


