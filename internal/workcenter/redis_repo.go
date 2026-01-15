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

func (r *RedisRepo) Set(ctx context.Context, id string, value models.WorkcenterDTO) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	pipe := r.client.Pipeline()
	pipe.Set(ctx, "workcenter:"+id, data, 0)
	pipe.SAdd(ctx, "workcenters", id)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.WorkcenterDTO, error) {
	data, err := r.client.Get(ctx, "workcenter:"+id).Bytes()
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
		// Even if some keys are missing (shouldn't happen with SADD/SREM consistency), we continue
		// But if pipe fails completely, return error
		return nil, err
	}

	var workcenters []models.WorkcenterDTO
	for _, cmd := range cmdMap {
		data, err := cmd.Bytes()
		if err == redis.Nil {
			continue // Should have been removed from Set, but skip for robustness
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

func (r *RedisRepo) Delete(ctx context.Context, id string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, "workcenter:"+id)
	pipe.SRem(ctx, "workcenters", id)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}