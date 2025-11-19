package status

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Repository struct {
	state     *state.State
	client    *redis.Client
	redisRepo *RedisRepo
	memoryRepo *MemoryRepo
}

func NewStatusRepository(state *state.State, client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	memoryRepo := NewMemoryRepository(state)
	return &Repository{
		state: state,
		client: client,
		redisRepo: redisRepo,
		memoryRepo: memoryRepo,
	}
}

func (r *Repository) Set(ctx context.Context, id string, value models.StatusDTO) error {
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

func (r *Repository) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error {
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

func (r *Repository) FindByID(ctx context.Context, id string) (models.StatusDTO, models.DataSource, error) {
	status, err := r.memoryRepo.FindByID(ctx, id)
	if err == nil {
		return status, models.SourceMemory, nil
	}
	if err != ErrStatusNotFound {
		return models.StatusDTO{}, models.SourceNone, err
	}

	status, err = r.redisRepo.FindByID(ctx, id)
	if err == nil {
		_ = r.memoryRepo.Set(ctx, id, status)
		return status, models.SourceRedis, nil
	}
	if err != ErrStatusNotFound {
		return models.StatusDTO{}, models.SourceNone, err
	}
	return models.StatusDTO{}, models.SourceNone, ErrStatusNotFound
}


