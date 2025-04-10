package services

import (
	"actions-service/internal"
	"actions-service/internal/clients"
	"sync"
)

type OperatorService interface{
	
}

type operatorService struct {
	client clients.HttpBackendClient
	state *internal.ServiceState
	mu *sync.Mutex
}

func NewOperatorService(client clients.HttpBackendClient, state *internal.ServiceState) StatusService{
	return &statusService{
		client: client,
		state: state,
		mu: &sync.Mutex{},
	}
}

