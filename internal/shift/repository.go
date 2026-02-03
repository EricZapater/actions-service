package shift

import (
	"actions-service/internal/models"
	"context"
	"sync"
	"time"
)

type Repository struct {	
	mu sync.RWMutex
	shifts map[string]models.ShiftDTO
}

func NewShiftRepository() *Repository {
	return &Repository{
		shifts: make(map[string]models.ShiftDTO),
	}
}

func (r *Repository) Set(ctx context.Context, id string, value models.ShiftDTO) {	
	r.mu.Lock()	
	defer r.mu.Unlock()
	r.shifts[id] = value
}

func (r *Repository) FindCurrent(ctx context.Context, now time.Time, id string) (models.ShiftDetailDTO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	shift, exists := r.shifts[id]
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


func (r *Repository) FindByID(ctx context.Context, id string) (models.ShiftDTO, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    shift, exists := r.shifts[id]
    if !exists {
        return models.ShiftDTO{}, ErrShiftNotFound
    }
    return shift, nil
}

func (r *Repository) FindShiftDetailByID(ctx context.Context, shiftID, detailID string) (models.ShiftDetailDTO, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    shift, exists := r.shifts[shiftID]
    if !exists {
        return models.ShiftDetailDTO{}, ErrShiftNotFound
    }
    
    for _, detail := range shift.ShiftDetails {
        if detail.ID.String() == detailID {
            return detail, nil
        }
    }
    
    return models.ShiftDetailDTO{}, ErrShiftDetailNotFound 
}