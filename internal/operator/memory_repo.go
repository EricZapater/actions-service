package operator

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


func(r *MemoryRepo) Set(ctx context.Context, id string, value models.OperatorDTO) error{
	r.state.Mu.Lock()			
	defer r.state.Mu.Unlock()
	r.state.Operators[id] = value
	return nil
}

func(r *MemoryRepo) FindByID(ctx context.Context, id string) (models.OperatorDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	operator, exists := r.state.Operators[id]
	if !exists {
		return models.OperatorDTO{}, ErrOperatorNotFound
	}
	return operator, nil
}

func(r *MemoryRepo) List(ctx context.Context) ([]models.OperatorDTO, error){
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	var operators []models.OperatorDTO
	for _, operator := range r.state.Operators {
		operators = append(operators, operator)
	}
	return operators, nil
}