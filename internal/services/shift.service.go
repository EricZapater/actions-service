package services

import (
	"actions-service/internal"
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
	GetShiftDetailsByShiftId(ctx context.Context, id uuid.UUID)([]models.ShiftDetail, error)
	GetDetailByIdBetweenHours(ctx context.Context, shiftId uuid.UUID, currentTime time.Time)(models.ShiftDetail, error)
	BuildShiftsDTO(ctx context.Context)(map[uuid.UUID]*models.Shift, error)
	FindShiftDetailinDTO(ctx context.Context, now time.Time, shiftId uuid.UUID, state *internal.ServiceState)(*models.ShiftDetail, error)
	
}

type shiftService struct {
	client clients.HttpBackendClient
}

func NewShiftService(client clients.HttpBackendClient) ShiftService{
	return &shiftService{client: client}
}

func (s *shiftService) GetShifts(ctx context.Context)([]models.Shift, error) {	
	url := fmt.Sprintf("/api/Shift")
	response, err := s.client.DoGetRequest(ctx, "GET", url)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return nil, err
	}
	defer func() { 
		if response.Body != nil {
			_ = response.Body.Close() 
		}
	}()
	if response.StatusCode == 200 {		
		var shifts []models.Shift
		if err := json.NewDecoder(response.Body).Decode(&shifts); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return shifts, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.Status)
	}
	
}
func(s *shiftService)GetShiftDetailsByShiftId(ctx context.Context, id uuid.UUID)([]models.ShiftDetail, error){
	url	:= fmt.Sprintf("/api/Shift/Detail/%v", id)
	response, err := s.client.DoGetRequest(ctx, "GET", url)	
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return nil, err
	}
	defer func() { 
		if response.Body != nil {
			_ = response.Body.Close() 
		}
	}()
	if response.StatusCode == 200 {		
		var shifts []models.ShiftDetail
		if err := json.NewDecoder(response.Body).Decode(&shifts); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return shifts, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.Status)
	}
}



func (s *shiftService) GetDetailByIdBetweenHours(ctx context.Context,shiftId uuid.UUID, currentTime time.Time)(models.ShiftDetail, error){
	log.Printf("GetDetailByIdBetweenHours: %v", shiftId)
	cTime := utils.StringAsAPointer(currentTime.Format("15:04:05"))	
	type Params struct{
		Id *uuid.UUID `json:"id"`
		CurrentTime string `json:"currentTime"`
	}
	params := Params{
		Id:          (&shiftId),
        CurrentTime: *cTime,
	}
	log.Printf("Params: %v", params)
	url := fmt.Sprintf("/api/Shift/Detail/ByIdBetweenHours?id=%s&currentTime=%s", shiftId.String(), *cTime)
	fmt.Println(params)
	var detail models.ShiftDetail
	response, err := s.client.DoPostRequest(ctx, "POST", url, params)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return detail, err
	}
	defer func() { 
		if response.Body != nil {
			_ = response.Body.Close() 
		}
	}()
	if response.StatusCode == 200 {		
		
		if err := json.NewDecoder(response.Body).Decode(&detail); err != nil {
			return detail, fmt.Errorf("error deserializing response: %w", err)
		}
		return detail, nil
	}else{
		return detail, fmt.Errorf("Response error: %v", response.Status)
	}
}

func(s *shiftService) BuildShiftsDTO(ctx context.Context)(map[uuid.UUID]*models.Shift, error){
	//shiftService := NewShiftService(s.client)
	ShiftDTO := make(map[uuid.UUID]*models.Shift)
	shifts, err := s.GetShifts(ctx)
	if err != nil {
		return nil, err
	}
	for _, shift := range shifts {
		if shift.Id == uuid.Nil {
			continue
		}
		//Recuperar detalls
		details, err := s.GetShiftDetailsByShiftId(ctx, shift.Id)
		if err != nil {
			return nil, err
		}
		//Construir el DTO
		ShiftDTO[shift.Id] = &models.Shift{
			Id: shift.Id,
			Name: shift.Name,
			ShiftDetail: details,
		}
	}
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
	state.Shifts = ShiftDTO
	return ShiftDTO, nil
}
func (s *shiftService)FindShiftDetailinDTO(ctx context.Context, now time.Time, shiftId uuid.UUID, state *internal.ServiceState)(*models.ShiftDetail, error){
	/*state := internal.GetInstance()
	state.Mu.Lock()	
	defer state.Mu.Unlock()*/
	var shiftDetail models.ShiftDetail
	shift, exists := state.Shifts[shiftId]
	if !exists {
		return &shiftDetail, fmt.Errorf("shift not found")
	}
	normalizedTime := time.Date(0,1,1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	for _, detail := range shift.ShiftDetail {
		
		startTime := time.Date(0, 1, 1, detail.StartTime.Hour(), detail.StartTime.Minute(), detail.StartTime.Second(), 0, now.Location())
		endTime := time.Date(0, 1, 1, detail.EndTime.Hour(), detail.EndTime.Minute(), detail.EndTime.Second(), 0, now.Location())

		// Comprovem si el torn creua la mitjanit
		if endTime.Before(startTime) { // Torn que passa després de mitjanit
			if normalizedTime.After(startTime) || normalizedTime.Before(endTime) {
				return &detail, nil
			}
		} else { // Torn normal
			if normalizedTime.After(startTime) && normalizedTime.Before(endTime) {
				return &detail, nil
			}
		}
	}
	return &shiftDetail, fmt.Errorf("shift not found")
}