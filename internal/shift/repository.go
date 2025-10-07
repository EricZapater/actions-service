package shift

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"time"
)

type Repository struct {	
	state *state.State
}

func NewShiftRepository(state *state.State) *Repository {
	return &Repository{
		state: state,
	}
}

func (r *Repository) Set(ctx context.Context, id string, value models.ShiftDTO) {	
	r.state.Mu.Lock()	
	defer r.state.Mu.Unlock()
	r.state.Shifts[id] = value
}

func (r *Repository) FindCurrent(ctx context.Context, now time.Time, id string) (models.ShiftDetailDTO, error) {
	r.state.Mu.RLock()
	defer r.state.Mu.RUnlock()
	shift, exists := r.state.Shifts[id]
	if !exists {
		return models.ShiftDetailDTO{}, ErrShiftNotFound
	}
	normalizedTime := time.Date(0,1,1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	for _, detail := range shift.ShiftDetails {
		startTime := time.Date(0, 1, 1, detail.StartTime.Hour(), detail.StartTime.Minute(), detail.StartTime.Second(), 0, now.Location())
		endTime := time.Date(0, 1, 1, detail.EndTime.Hour(), detail.EndTime.Minute(), detail.EndTime.Second(), 0, now.Location())

		if endTime.Before(startTime) { 
			if normalizedTime.After(startTime) || normalizedTime.Before(endTime) {
				return detail, nil
			}
		} else {
			if normalizedTime.After(startTime) && normalizedTime.Before(endTime) {
				return detail, nil
			}
		}
	}
	return models.ShiftDetailDTO{}, ErrShiftNotFound
}

