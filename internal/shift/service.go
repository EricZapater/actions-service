package shift

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	BuildDTO(ctx context.Context) error
	getShiftDetails(ctx context.Context, shiftID uuid.UUID) ([]models.ShiftDetailDTO, error)
	FindCurrentShift(ctx context.Context, now time.Time, shiftID uuid.UUID) (models.ShiftDetailDTO, error)
}

type service struct {
	client clients.HttpBackendClient
	repo  Repository
}

func NewShiftService(client clients.HttpBackendClient, repo Repository) Service {
	return &service{
		client: client,
		repo:  repo,
	}
}

func (s *service) BuildDTO(ctx context.Context)error {	
	url := "/api/Shift"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		log.Printf("Something went wrong calling the backend %v", err)
		return err
	}
	defer func() { 
		if response.Body != nil {
			_ = response.Body.Close() 
		}
	}()
	var shifts []models.ShiftDTO
	if response.StatusCode == 200 {		
		
		if err := json.NewDecoder(response.Body).Decode(&shifts); err != nil {
			return fmt.Errorf("error deserializing response: %w", err)
		}		
	}else{
		return fmt.Errorf("response error: %v", response.Status)
	}

	for i, shift := range shifts {
		details, err := s.getShiftDetails(ctx, shift.ID)
		if err != nil {
			return fmt.Errorf("error getting shift details for shift %s: %w", shift.ID, err)
		}
		shifts[i].ShiftDetails = details
	}
	for _, shift := range shifts {
		s.repo.Set(ctx, shift.ID.String(), shift)
	}
	return nil	
}

func(s *service)getShiftDetails(ctx context.Context, id uuid.UUID)([]models.ShiftDetailDTO, error){
	url	:= fmt.Sprintf("/api/Shift/Detail/%v", id)
	response, err := s.client.DoGetRequest(ctx, url)	
	if err != nil {
		log.Printf("Something went wrong calling the backend %v", err)
		return nil, err
	}
	defer func() { 
		if response.Body != nil {
			_ = response.Body.Close() 
		}
	}()
	if response.StatusCode == 200 {		
		var shifts []models.ShiftDetailDTO
		if err := json.NewDecoder(response.Body).Decode(&shifts); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return shifts, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.Status)
	}
}

func (s *service) FindCurrentShift(ctx context.Context, now time.Time, shiftID uuid.UUID) (models.ShiftDetailDTO, error) {
	shiftDetail, err := s.repo.FindCurrent(ctx, now, shiftID.String())
	if err != nil {
		return models.ShiftDetailDTO{}, err
	}
	return shiftDetail, nil
}