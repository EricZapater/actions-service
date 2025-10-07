package state

import (
	"actions-service/internal/models"
	"sync"
)

type State struct {
	Mu sync.RWMutex
	Workcenters map[string]models.WorkcenterDTO
	Shifts map[string]models.ShiftDTO
	Operators map[string]models.OperatorDTO
}

func New() *State{
	return &State{
		Workcenters: make(map[string]models.WorkcenterDTO),
		Shifts: make(map[string]models.ShiftDTO),
		Operators: make(map[string]models.OperatorDTO),
	}
}