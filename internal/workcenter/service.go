package workcenter

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/shift"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	BuildDTO(ctx context.Context) error
	SetCurrentShift(ctx context.Context) error
}

type service struct {
	client clients.HttpBackendClient
	repo  Repository
	shiftService shift.Service
}

func NewWorkcenterService(client clients.HttpBackendClient, repo Repository, shiftService shift.Service) Service {
	return &service{
		client: client,
		repo:  repo,
		shiftService: shiftService,
	}
}

func (s *service) BuildDTO(ctx context.Context)error {
	url := "/api/Workcenter"
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
	var workcenters []models.Workcenter
	if response.StatusCode == 200 {
		if err := json.NewDecoder(response.Body).Decode(&workcenters); err != nil {
			return fmt.Errorf("error deserializing response: %w", err)
		}
	}else{
		return fmt.Errorf("response error: %v", response.Status)
	}
	
	//recuperar l'area?
	
	for _, workcenter := range workcenters {
		if workcenter.ShiftId == uuid.Nil {
			//esborrar-lo del redis i de la memoria
			continue
		}
		//montar el DTO
		WorkcentersDTO := &models.WorkcenterDTO{
			WorkcenterID: workcenter.Id,
			WorkcenterName: workcenter.Name,
			WorkcenterDescription: workcenter.Description,
			AreaID: workcenter.AreaId,
			AreaDescription: "",
			ShiftID: workcenter.ShiftId,
			ShiftName: "",
			ShiftDetailId: uuid.Nil,			
			StatusID: uuid.Nil,
			StatusName: "",
			StatusOperatorsAllowed: true,
			StatusClosed: true,
			StatusStopped: true,
			StatusColor: "",
			StatusStartTime: time.Now(),
		}
		wc, source, err := s.repo.FindByID(ctx, WorkcentersDTO.WorkcenterID.String())
		if err != nil {
			if err == ErrWorkcenterNotFound {
				//setejar-lo
				if err := s.repo.Set(ctx, WorkcentersDTO.WorkcenterID.String(), *WorkcentersDTO); err != nil {
					return fmt.Errorf("error setting workcenter %s: %w", WorkcentersDTO.WorkcenterID.String(), err)
				}
			}
		}
		//si no hi ha error, per tant existeix
		if source == models.SourceMemory {
			//comprovar areaid, shiftid
			if wc.AreaID != WorkcentersDTO.AreaID || wc.ShiftID != WorkcentersDTO.ShiftID {	
				//Ull!! aqui caldria fer una crida al back per fer el nou shift

				WorkcentersDTO.StatusID = wc.StatusID
				WorkcentersDTO.StatusName = wc.StatusName
				WorkcentersDTO.StatusOperatorsAllowed = wc.StatusOperatorsAllowed
				WorkcentersDTO.StatusClosed = wc.StatusClosed
				WorkcentersDTO.StatusStopped = wc.StatusStopped
				WorkcentersDTO.StatusColor = wc.StatusColor
				WorkcentersDTO.StatusStartTime = wc.StatusStartTime
				if err := s.repo.Set(ctx, WorkcentersDTO.WorkcenterID.String(), *WorkcentersDTO); err != nil {
					return fmt.Errorf("error updating workcenter %s: %w", WorkcentersDTO.WorkcenterID.String(), err)
				}
			}
		}
		if source == models.SourceRedis {
			if err := s.repo.Set(ctx, WorkcentersDTO.WorkcenterID.String(), wc); err != nil {
				return fmt.Errorf("error updating workcenter %s: %w", WorkcentersDTO.WorkcenterID.String(), err)
			}
		}

	//tanca el bucle	
	}
	wc, err := s.repo.List(ctx) 
	if err != nil {
		return fmt.Errorf("error listing workcenters: %w", err)
	}
	log.Printf("Loaded %d workcenters from backend", len(wc))
	fmt.Println(wc)
	return nil
}

func(s *service) SetCurrentShift(ctx context.Context)error{
	workcenters, err := s.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing workcenters: %w", err)
	}
	for _, wc := range workcenters {		
		now := time.Now()
		shiftDetail, err := s.shiftService.FindCurrentShift(ctx, now, wc.ShiftID)
		if err != nil {
			return fmt.Errorf("error finding current shift for workcenter %s: %w", wc.WorkcenterID.String(), err)
		}
		if wc.ShiftDetailId != shiftDetail.ID {
			start := time.Date(
				now.Year(), now.Month(), now.Day(),
				shiftDetail.StartTime.Hour(), shiftDetail.StartTime.Minute(), shiftDetail.StartTime.Second(),
				0, now.Location(),
			)
			wc.ShiftDetailId = shiftDetail.ID
			wc.ShiftDetailStartTime = shiftDetail.StartTime
			wc.ShiftDetailEndTime = shiftDetail.EndTime
			wc.ShiftDetailIsProductiveTime = shiftDetail.IsProductiveTime
			//backend
			request := models.CreateWorkcenterShiftDTO{
				WorkcenterID:  wc.WorkcenterID,
				ShiftDetailId: shiftDetail.ID,
				StartTime:     start.Format("2006-01-02T15:04:05"),
			}
			if err := s.createWorkcenterShift(ctx, request); err != nil {
				return fmt.Errorf("error creating workcenter shift for workcenter %s: %w", wc.WorkcenterID.String(), err)
			}
			//memory
			if err := s.repo.Set(ctx, wc.WorkcenterID.String(), wc); err != nil {
				return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
			}
		}
	}
	return nil
}

func(s *service) createWorkcenterShift(ctx context.Context, request models.CreateWorkcenterShiftDTO)error{
	requests := []models.CreateWorkcenterShiftDTO{request}
	//log.Printf("%v",requests)
	start := time.Now()
	_, err := s.client.DoPostRequest(ctx, "/api/WorkcenterShift/CreateWorkcenterShifts", requests)
	if err != nil {
		log.Printf("Something went wrong calling the backend %v", err)
		return err
	}
	end := time.Now()
	log.Printf("CreateWorkcenterShift took %v", end.Sub(start))
	//log.Printf("Response status: %v\n", response.Status)
	return nil
}


func(s *service) GetAreas(ctx context.Context)([]models.Area, error){		
	response, err := s.client.DoGetRequest(ctx, "/api/area")
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
		var areas []models.Area
		if err := json.NewDecoder(response.Body).Decode(&areas); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return areas, nil
	}else{
		return nil, fmt.Errorf("Response error: %v", response.Status)
	}
	
}