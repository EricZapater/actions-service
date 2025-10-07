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

func(r *RedisRepo) Set(ctx context.Context, id string, value models.OperatorDTO) error{
	data, err := json.Marshal(value)	
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, id, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id string) (models.OperatorDTO, error){
	data, err := r.client.Get(ctx, id).Bytes()
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