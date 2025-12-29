package workorderphase

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

func(r *MemoryRepo) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error{
	r.state.Mu.Lock()			
	defer r.state.Mu.Unlock()
	r.state.Workcenters[id] = value
	return nil
}