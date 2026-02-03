package workorderphase

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

func (r *RedisRepo) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
    data, err := json.Marshal(value)
	fmt.Println(&value)
    if err != nil {
        return err
    }
    if err := r.client.Set(ctx, "workcenter:"+id, data, 0).Err(); err != nil {
        return err
    }
    return nil
}

func (r *RedisRepo) List(ctx context.Context) ([]models.WorkcenterDTO, error) {
	// 1. Get all IDs from the Set
	ids, err := r.client.SMembers(ctx, "workcenters").Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []models.WorkcenterDTO{}, nil
	}

	// 2. Fetch all workcenters in a pipeline
	pipe := r.client.Pipeline()
	cmdMap := make(map[string]*redis.StringCmd)
	
	for _, id := range ids {
		cmdMap[id] = pipe.Get(ctx, "workcenter:"+id)
	}
	
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	var workcenters []models.WorkcenterDTO
	for _, cmd := range cmdMap {
		data, err := cmd.Bytes()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err 
		}
		
		var wc models.WorkcenterDTO
		if err := json.Unmarshal(data, &wc); err == nil {
			workcenters = append(workcenters, wc)
		}
	}
	return workcenters, nil
}


