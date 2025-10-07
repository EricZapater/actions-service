package workcenter

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

func NewWorkcenterRepository(state *state.State, client *redis.Client) *Repository {
	redisRepo := NewRedisRepository(client)
	memoryRepo := NewMemoryRepository(state)
	return &Repository{
		state: state,
		client: client,
		redisRepo: redisRepo,
		memoryRepo: memoryRepo,
	}
}

func(r *Repository) Set(ctx context.Context, id string, value models.WorkcenterDTO)error{
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

func (r *Repository) FindByID(ctx context.Context, id string) (models.WorkcenterDTO, models.DataSource, error) {
	workcenter, err := r.memoryRepo.FindByID(ctx, id)
	if err == nil {
		return workcenter, models.SourceMemory, nil
	}
	if err != ErrWorkcenterNotFound {
		return models.WorkcenterDTO{}, models.SourceNone, err
	}

	workcenter, err = r.redisRepo.FindByID(ctx, id)
	if err == nil {
		_ = r.memoryRepo.Set(ctx, id, workcenter)
		return workcenter, models.SourceRedis, nil
	}
	if err != ErrWorkcenterNotFound {
		return models.WorkcenterDTO{}, models.SourceNone, err
	}
	return models.WorkcenterDTO{}, models.SourceNone, ErrWorkcenterNotFound
}


func(r *Repository) List(ctx context.Context) ([]models.WorkcenterDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	var workcenters []models.WorkcenterDTO
	for _, workcenter := range r.state.Workcenters {
		workcenters = append(workcenters, workcenter)
	}
	return workcenters, nil
}