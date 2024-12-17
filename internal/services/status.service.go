package services

import (
	"actions-service/internal"
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type StatusService interface {
	GetStatusById(ctx context.Context, statusID uuid.UUID)(models.Status, error)
	UpdateWorkcenterStatus(ctx context.Context, workcenterID, statusID uuid.UUID) error
}

type statusService struct {
	client *clients.ClientWithResponses
	state *internal.ServiceState
	mu *sync.Mutex
}

func NewStatusService(client *clients.ClientWithResponses, state *internal.ServiceState) StatusService{
	return &statusService{
		client: client,
		state: state,
		mu: &sync.Mutex{},
	}
}

func(s *statusService)UpdateWorkcenterStatus(ctx context.Context, workcenterID, statusID uuid.UUID) error {
	/*s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(s.state.Workcenters)
	wc, exists := s.state.Workcenters[workcenterID]
	fmt.Println("1",wc)*/
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
	
	wc, exists := state.Workcenters[workcenterID]
	
	if !exists{
		return fmt.Errorf("workcenter with ID %v not found", workcenterID)
	}
	status, err := s.GetStatusById(ctx, statusID)
	if err != nil{
		return err
	}
	wc.StatusId = status.Id
	wc.StatusName = status.Name
	wc.StatusOperatorsAllowed = status.OperatorsAllowed
	wc.StatusClosed = status.Closed
	wc.StatusStopped = status.Stopped
	wc.StatusColor = status.Color	
	wc.StatusStartTime = time.Now()
	return nil
}

func(s *statusService)GetStatusById(ctx context.Context, statusID uuid.UUID)(models.Status, error){
	response, err := s.client.GetApiMachineStatusIdWithResponse(ctx, statusID)
	var status models.Status
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return status, err
	}
	defer func() { 
		if response.HTTPResponse.Body != nil {
			_ = response.HTTPResponse.Body.Close() 
		}
	}()
	if response.HTTPResponse.StatusCode == 200 {		
		
		if err := json.Unmarshal(response.Body, &status); err != nil {
			return status, fmt.Errorf("error deserializing response: %w", err)
		}
		return status, nil
	}else{
		return status, fmt.Errorf("response error: %v", response.HTTPResponse.Status)
	}
}

