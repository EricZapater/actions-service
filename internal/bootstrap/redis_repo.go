package bootstrap

import (
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"

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

// FlushDB clears all data from the current Redis database
func (r *RedisRepo) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
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

func (r *RedisRepo) SetMultiple(ctx context.Context, workcenters []models.WorkcenterDTO) error {
    pipe := r.client.Pipeline()
    
    for _, wc := range workcenters {
        data, err := json.Marshal(wc)
        if err != nil {
            return err
        }
        pipe.Set(ctx, fmt.Sprintf("workcenter:%s", wc.WorkcenterID), data, 0)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}