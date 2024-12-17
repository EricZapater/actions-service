package services

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ShiftService interface {
	GetShifts(ctx context.Context)([]models.Shift, error)
	GetDetailByIdBetweenHours(ctx context.Context, shiftId uuid.UUID, currentTime time.Time)(models.ShiftDetail, error)
	
}

type shiftService struct {
	client *clients.ClientWithResponses
}

func NewShiftService(client *clients.ClientWithResponses) ShiftService{
	return &shiftService{client: client}
}

func (s *shiftService) GetShifts(ctx context.Context)([]models.Shift, error) {	
	response, err := s.client.GetApiShiftWithResponse(ctx)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return nil, err
	}
	defer func() { 
		if response.HTTPResponse.Body != nil {
			_ = response.HTTPResponse.Body.Close() 
		}
	}()
	if response.HTTPResponse.StatusCode == 200 {		
		var shifts []models.Shift
		if err := json.Unmarshal(response.Body, &shifts); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return shifts, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.HTTPResponse.Status)
	}
	
}



func (s *shiftService) GetDetailByIdBetweenHours(ctx context.Context,shiftId uuid.UUID, currentTime time.Time)(models.ShiftDetail, error){
	cTime := utils.StringAsAPointer(currentTime.Format("15:04:05"))	
	params := &clients.PostApiShiftDetailByIdBetweenHoursParams{
		Id:          (&shiftId),
        CurrentTime: cTime,
	}
	var detail models.ShiftDetail
	response, err := s.client.PostApiShiftDetailByIdBetweenHoursWithResponse(ctx, params)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return detail, err
	}
	defer func() { 
		if response.HTTPResponse.Body != nil {
			_ = response.HTTPResponse.Body.Close() 
		}
	}()
	if response.HTTPResponse.StatusCode == 200 {		
		
		if err := json.Unmarshal(response.Body, &detail); err != nil {
			return detail, fmt.Errorf("error deserializing response: %w", err)
		}
		return detail, nil
	}else{
		return detail, fmt.Errorf("Response error: %v", response.HTTPResponse.Status)
	}
}