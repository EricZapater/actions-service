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
	CheckWorkcenterShift(ctx context.Context)error
}

type workcenterService struct {
	client clients.HttpBackendClient
}

func NewWorkcenterService(client clients.HttpBackendClient) WorkcenterService{
	return &workcenterService{client: client}
}

func (s *workcenterService) GetWorkcenters(ctx context.Context)([]models.Workcenter, error){		
	url := fmt.Sprintf("/api/Workcenter")
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
		var workcenters []models.Workcenter
		if err := json.NewDecoder(response.Body).Decode(&workcenters); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return workcenters, nil
	}else{
		return nil, fmt.Errorf("response error: %v", response.Status)
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
		now := time.Now()
		start := time.Date(
			now.Year(), now.Month(), now.Day(),
			detail.StartTime.Hour(), detail.StartTime.Minute(), detail.StartTime.Second(),
			0, now.Location(),
		)
		request := models.CreateWorkcenterShiftDTO{
			WorkcenterID: wc.Id,
			ShiftDetailId: detail.Id,
			StartTime: start.Format("2006-01-02T15:04:05"),
		}
		err = s.CreateWorkcenterShift(ctx, request)
			if err != nil {
				log.Fatalf("error creating workcenter shift %v", err)
				return nil, err
			}
		shift := cache.Shifts[wc.ShiftId]
		
		area := cache.Areas[wc.AreaId]
				
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
	}
	state := internal.GetInstance() // Obtenim el singleton
	state.Mu.Lock()
	defer state.Mu.Unlock()
    state.Workcenters = WorkcentersDTO
	return WorkcentersDTO, nil
}

func (s *workcenterService) CheckWorkcenterShift(ctx context.Context)error{
	startProcessTime := time.Now()
	fmt.Printf("Inici process %v\n", startProcessTime)
	var endProcessTime time.Time
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
	now := time.Now()
	shiftService := NewShiftService(s.client)
	//Comparar per cada maquina el torn del que ve, i canviar-li si cal
	for _, wc := range state.Workcenters {
		startTime := time.Date(0,1,1,wc.ShiftDetailStartTime.Hour(), wc.ShiftDetailStartTime.Minute(), wc.ShiftDetailStartTime.Second(), 0, now.Location())
		endTime := time.Date(0,1,1, wc.ShiftDetailEndTime.Hour(), wc.ShiftDetailEndTime.Minute(), wc.ShiftDetailEndTime.Second(), 0, now.Location())
		startDate := time.Date(now.Year(), now.Month(), now.Day(), wc.ShiftDetailStartTime.Hour(), wc.ShiftDetailStartTime.Minute(), wc.ShiftDetailStartTime.Second(), 0, now.Location())
		var endDate time.Time
		if endTime.Before(startTime){
			nextNow := now.Add(24*time.Hour)
			endDate = time.Date(nextNow.Year(), nextNow.Month(), nextNow.Day(), wc.ShiftDetailEndTime.Hour(), wc.ShiftDetailEndTime.Minute(), wc.ShiftDetailEndTime.Second(), 0, now.Location())
		}else{
			endDate = time.Date(now.Year(), now.Month(), now.Day(), wc.ShiftDetailEndTime.Hour(), wc.ShiftDetailEndTime.Minute(), wc.ShiftDetailEndTime.Second(), 0, now.Location())
		}
		if now.After(startDate) && now.Before(endDate){
			fmt.Println("dins el torn")
		}else{
			fmt.Println("fora del torn")
			//canviar el torn
			//trobar el torn que toca
			detail, err := shiftService.FindShiftDetailinDTO(ctx, now, wc.ShiftId, state)
			if err != nil {
				log.Fatalf("error checking details in DTO %v", err)
				return err				
			}
			//actualitzar el torn al back
			start := time.Date(
				now.Year(), now.Month(), now.Day(),
				detail.StartTime.Hour(), detail.StartTime.Minute(), detail.StartTime.Second(),
				0, now.Location(),
			)
			request := models.CreateWorkcenterShiftDTO{
				WorkcenterID: wc.WorkcenterId,
				ShiftDetailId: detail.Id,
				StartTime:  start.Format("2006-01-02T15:04:05"),
			}
			err = s.CreateWorkcenterShift(ctx, request)
			if err != nil {
				log.Fatalf("error creating workcenter shift %v", err)
				return err
			}
			//actualitzar el torn al DTO
			wc.ShiftDetailId = detail.Id
			wc.ShiftDetailStartTime = detail.StartTime
			wc.ShiftDetailEndTime = detail.EndTime
			wc.ShiftDetailIsProductiveTime = detail.IsProductiveTime			
		}
		fmt.Printf("Workcenter %s, shiftname: %s, starttime: %v, endtime %v\n", wc.WorkcenterDescription, wc.ShiftName, wc.ShiftDetailStartTime, wc.ShiftDetailEndTime)
	}
	endProcessTime = time.Now()
	fmt.Println("checkshift duration: ",endProcessTime.Sub(startProcessTime))	
	return nil
}

func (s *workcenterService) CreateWorkcenterShift(ctx context.Context, request models.CreateWorkcenterShiftDTO)error{
	requests := []models.CreateWorkcenterShiftDTO{request}
	log.Printf("%v",requests)
	response, err := s.client.DoPostRequest(ctx, "POST", "/api/WorkcenterShift/CreateWorkcenterShifts", requests)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
		return err
	}
	log.Printf("Response status: %v\n", response.Status)
	return nil
}