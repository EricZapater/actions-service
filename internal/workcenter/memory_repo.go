package workcenter

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
)

type MemoryRepo struct {
	state *state.State
}

func NewMemoryRepository(state *state.State) *MemoryRepo {
	return &MemoryRepo{
		state: state,
	}
}

func(r *MemoryRepo) Set(ctx context.Context, id string, value models.WorkcenterDTO) error{
	r.state.Mu.Lock()	
	defer r.state.Mu.Unlock()
	r.state.Workcenters[id] = value
	return nil
}

func(r *MemoryRepo) FindByID(ctx context.Context, id string) (models.WorkcenterDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	workcenter, exists := r.state.Workcenters[id]
	if !exists {
		return models.WorkcenterDTO{}, ErrWorkcenterNotFound
	}
	return workcenter, nil
}