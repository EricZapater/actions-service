package services

import (
	"actions-service/internal"
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type cache struct {
	Shifts map[uuid.UUID]models.Shift	
	Areas map[uuid.UUID]models.Area
}
type WorkcenterService interface {
	GetWorkcenters(ctx context.Context) ([]models.Workcenter, error)
	BuildWorkcenterDTO(ctx context.Context)(map[uuid.UUID]*models.WorkcenterDTO, error)
}

type workcenterService struct {
	client *clients.ClientWithResponses
}

func NewWorkcenterService(client *clients.ClientWithResponses) WorkcenterService{
	return &workcenterService{client: client}
}

func (s *workcenterService) GetWorkcenters(ctx context.Context)([]models.Workcenter, error){		
	response, err := s.client.GetApiWorkcenterWithResponse(ctx)
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
		var workcenters []models.Workcenter
		if err := json.Unmarshal(response.Body, &workcenters); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return workcenters, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.HTTPResponse.Status)
	}
}

func (s *workcenterService) BuildWorkcenterDTO(ctx context.Context)(map[uuid.UUID]*models.WorkcenterDTO, error) {
	shiftService := NewShiftService(s.client)
	areaService := NewAreaService(s.client)
	cache := cache{
		Shifts: make(map[uuid.UUID]models.Shift),		
		Areas: make(map[uuid.UUID]models.Area),
	}
	fmt.Println("Getting shifts")
	shifts, err := shiftService.GetShifts(ctx)
	if err != nil{
		return nil, err
	}
	for _, shift := range shifts {
		cache.Shifts[shift.Id] = shift
	}
	fmt.Println("Getting areas")
	areas, err := areaService.GetAreas(ctx)
	if err != nil {
		return nil, err
	}
	for _, area := range areas {
		cache.Areas[area.Id] = area
	}
	fmt.Println("Getting workcenters")
	workcenters, err := s.GetWorkcenters(ctx)
	if err != nil {
		return nil, err
	}
	
	WorkcentersDTO := make(map[uuid.UUID]*models.WorkcenterDTO)
	for _, wc := range workcenters {
		if wc.ShiftId == uuid.Nil || wc.Disabled{
			continue
		}
		detail, err  := shiftService.GetDetailByIdBetweenHours(ctx, wc.ShiftId, time.Now())
		if err != nil {
			log.Printf("error retrieving details: %v\n", err)
		}
		shift := cache.Shifts[wc.ShiftId]
		
		area := cache.Areas[wc.AreaId]
		
		fmt.Printf("Constructing DTO for workcenter %s\n", wc.Description)
		WorkcentersDTO[wc.Id] = &models.WorkcenterDTO{
			WorkcenterId: wc.Id,
			WorkcenterName: wc.Name,
			WorkcenterDescription: wc.Description,
			AreaId: area.Id,
			AreaDescription: area.Description,
			ShiftId: shift.Id,
			ShiftName: shift.Name,
			ShiftDetailId: detail.Id,
			ShiftDetailStartTime: detail.StartTime,
			ShiftDetailEndTime: detail.EndTime,
			ShiftDetailIsProductiveTime: detail.IsProductiveTime,
			StatusId: uuid.Nil,
			StatusName: "",
			StatusOperatorsAllowed: true,
			StatusClosed: true,
			StatusStopped: true,
			StatusColor: "",
			StatusStartTime: time.Now(),
		}
		fmt.Printf("DTO created: %+v\n", WorkcentersDTO[wc.Id])
	}
	state := internal.GetInstance() // Obtenim el singleton
	state.Mu.Lock()
	defer state.Mu.Unlock()
    state.Workcenters = WorkcentersDTO
	return WorkcentersDTO, nil
}