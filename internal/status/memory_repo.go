package status

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

func (r *MemoryRepo) Set(ctx context.Context, id string, value models.StatusDTO)  error{
	r.state.Mu.Lock()
	defer r.state.Mu.Unlock()
	r.state.Statuses[id] = value
	return nil
}


func(r *MemoryRepo) SetWorkcenterDTO(ctx context.Context, id string, value models.WorkcenterDTO) error{
	r.state.Mu.Lock()			
	defer r.state.Mu.Unlock()
	r.state.Workcenters[id] = value
	return nil
}

func(r *MemoryRepo) FindByID(ctx context.Context, id string) (models.StatusDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	status, exists := r.state.Statuses[id]
	if !exists {
		return models.StatusDTO{}, ErrStatusNotFound
	}
	return status, nil
}

func(r *MemoryRepo) List(ctx context.Context)([]models.StatusDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	var statuses []models.StatusDTO
	for _, status := range r.state.Statuses{
		statuses = append(statuses, status)
	}
	return statuses, nil
}