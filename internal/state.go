package internal

import (
	"actions-service/internal/models"
	"sync"

	"github.com/google/uuid"
)

var instance *ServiceState
var once sync.Once

type ServiceState struct {
	Workcenters map[uuid.UUID]*models.WorkcenterDTO
	Mu sync.RWMutex
}

func GetInstance()*ServiceState {
	once.Do(func(){
		instance =&ServiceState{
			Workcenters: make(map[uuid.UUID]*models.WorkcenterDTO),
		}
	})
	return instance
}