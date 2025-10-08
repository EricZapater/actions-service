package operator

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Repository struct {
	state *state.State
	client *redis.Client
	redisRepo *RedisRepo
	memoryRepo *MemoryRepo
}

func NewOperatorRepository(state *state.State, client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	memoryRepo := NewMemoryRepository(state)
	return &Repository{
		state: state,
		client: client,
		redisRepo: redisRepo,
		memoryRepo: memoryRepo,
	}
}

func(r *Repository) Set(ctx context.Context, id string, value models.OperatorDTO)error{
	g, ctx := errgroup.WithContext(ctx)
	
	g.Go(func() error {
		return r.memoryRepo.Set(ctx, id, value)
		
	})

	g.Go(func() error {
		return r.redisRepo.Set(ctx, id, value)
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func(r *Repository) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO)error{
	g, ctx := errgroup.WithContext(ctx)
	
	g.Go(func() error {
		return r.memoryRepo.SetWorkcenterDTO(ctx, id, value)
		
	})

	g.Go(func() error {
		return r.redisRepo.SetWorkcenterDTO(ctx, id, value)
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (models.OperatorDTO, models.DataSource, error){
	operator, err := r.memoryRepo.FindByID(ctx, id)
	if err == nil {
		return operator, models.SourceMemory, nil
	}
	if err != ErrOperatorNotFound {
		return models.OperatorDTO{}, models.SourceNone, err
	}

	operator, err = r.redisRepo.FindByID(ctx, id)
	if err == nil {
		_ = r.memoryRepo.Set(ctx, id, operator)
		return operator, models.SourceRedis, nil
	}
	if err != ErrOperatorNotFound {
		return models.OperatorDTO{}, models.SourceNone, err
	}
	return models.OperatorDTO{}, models.SourceNone, ErrOperatorNotFound
}