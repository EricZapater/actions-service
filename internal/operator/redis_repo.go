package operator

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

func (r *RedisRepo) Set(ctx context.Context, id string, value models.OperatorDTO) error {
	data, err := json.Marshal(value)	
	if err != nil {
		return err
	}
	
	pipe := r.client.Pipeline()
	pipe.Set(ctx, "operator:"+id, data, 0)
	pipe.SAdd(ctx, "operators", id)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
	// TODO: This method seems misplaced or redundant as it deals with Workcenters in Operator repo.
	// Applying namespacing for consistency if it remains used.
	data, err := json.Marshal(value)	
	if err != nil {
		return err
	}
	// Assuming this was meant to update a workcenter from operator context?
	// Using consistent workcenter namespace just in case.
	if err := r.client.Set(ctx, "workcenter:"+id, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.OperatorDTO, error) {
	data, err := r.client.Get(ctx, "operator:"+id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.OperatorDTO{}, ErrOperatorNotFound
		}
		return models.OperatorDTO{}, err
	}

	var operator models.OperatorDTO
	if err := json.Unmarshal(data, &operator); err != nil {
		return models.OperatorDTO{}, err
	}
	return operator, nil
}