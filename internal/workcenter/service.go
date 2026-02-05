package workcenter

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/observability"
	"actions-service/internal/shift"
	"actions-service/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	BuildDTO(ctx context.Context) error
	SetCurrentShift(ctx context.Context) error
	GetWorkcenterDTO(ctx context.Context, id string) (*models.WorkcenterDTO, error)
	GetAllWorkcenters(ctx context.Context) ([]models.WorkcenterDTO, error)
}

// ... existing struct and constructor ...
type service struct {
	client clients.HttpBackendClient
	repo  Repository
	shiftService shift.Service
	statusPort StatusPort
	hub *ws.Hub
	logger *slog.Logger
	
}

func NewWorkcenterService(client clients.HttpBackendClient, repo Repository, shiftService shift.Service,statusPort StatusPort, hub *ws.Hub) Service {
	return &service{
		client: client,
		repo:  repo,
		shiftService: shiftService,
		statusPort: statusPort,
		hub: hub,
		logger: observability.NewLogger("info"),
	}
}

func (s *service) BuildDTO(ctx context.Context)error {
	// Get Backend workcenters
	url := "/api/workcenter"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		s.logger.ErrorContext(ctx,  "Failed to get workcenters from backend",
    	slog.String("error", err.Error()),
    	slog.String("url", url),)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		s.logger.ErrorContext(ctx,  "Failed to get workcenters from backend",
    	slog.String("error", response.Status),
    	slog.String("url", url),)
		return fmt.Errorf("error getting workcenters from backend: %s", response.Status)
	}

	var backendWorkcenters []models.Workcenter
	if err := json.NewDecoder(response.Body).Decode(&backendWorkcenters); err != nil {
		s.logger.ErrorContext(ctx,  "Failed to decode response from backend",
      	slog.String("error", err.Error()),
      	slog.String("url", url),)
		return fmt.Errorf("error decoding response: %w", err)
	}

	//Get Cache Workcenters
	cacheWorkcenters, err := s.repo.List(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx,  "Failed to get workcenters from cache",
     	slog.String("error", err.Error()),)
		return fmt.Errorf("error listing workcenters: %w", err)
	}

	//Map workcenters
	backendMap := make(map[string]models.Workcenter)	
	for _, wc := range backendWorkcenters {
		if wc.ShiftId == uuid.Nil {
			continue
		}
		backendMap[wc.Id.String()] = wc
	}
	cacheMap := make(map[string]models.WorkcenterDTO)
	for _, wc := range cacheWorkcenters {
		cacheMap[wc.WorkcenterID.String()] = wc
	}

	//Delete from cache                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             
	for cacheID := range cacheMap {
		if _, exists := backendMap[cacheID]; !exists {			
			if err := s.repo.Delete(ctx, cacheID); err != nil {
				s.logger.ErrorContext(ctx, "Failed to delete workcenter from cache",
					slog.String("workcenter_id", cacheID),
					slog.String("error", err.Error()),)
			}
			s.logger.InfoContext(ctx, "Workcenter deleted from cache",
				slog.String("workcenter_id", cacheID),)				
		}
	}

	//get default status
	defaultStatus, err := s.statusPort.GetDefaultStatus(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get default status",
			slog.String("error", err.Error()),)
		return fmt.Errorf("error getting default status: %w", err)
	}
	//Add to cache
	for backendID, backendWC := range backendMap {
		cacheWC, exists := cacheMap[backendID]
		//Add to cache
		if !exists {
			newWC := models.WorkcenterDTO{
				WorkcenterID: backendWC.Id,
				WorkcenterName: backendWC.Name,
				ShiftID: backendWC.ShiftId,
				ShiftDetailId: uuid.Nil,
				StatusID:              defaultStatus.StatusId,
                StatusName:            defaultStatus.Description,
                StatusOperatorsAllowed: defaultStatus.OperatorsAllowed,
                StatusClosed:          defaultStatus.Closed,
                StatusStopped:         defaultStatus.Stopped,
                StatusColor:           defaultStatus.Color,
                StatusStartTime:       time.Now(),				
			}
			/*parsedStatusID:= defaultStatus.StatusId.String()
			if err := s.statusPort.StatusIn(ctx, backendID, parsedStatusID, nil); err != nil {
				s.logger.ErrorContext(ctx, "Failed to set status in cache",
					slog.String("workcenter_id", backendID),
					slog.String("error", err.Error()),)
			}*/
			if err := s.repo.Set(ctx, backendID, newWC); err != nil {
				s.logger.ErrorContext(ctx, "Failed to add workcenter to cache",
					slog.String("workcenter_id", backendID),
					slog.String("error", err.Error()),)
			}
		}else{
			//Update cache
			needsUpdate := false
			updatedWC := cacheWC
			//Configuració
			if cacheWC.WorkcenterName != backendWC.Name {
				needsUpdate = true
				updatedWC.WorkcenterName = backendWC.Name
			}
			if cacheWC.WorkcenterDescription != backendWC.Description {
				needsUpdate = true
				updatedWC.WorkcenterDescription = backendWC.Description
			}
			if cacheWC.AreaID != backendWC.AreaId {
				needsUpdate = true
				updatedWC.AreaID = backendWC.AreaId
			}			
			if cacheWC.MultiWoAvailable != backendWC.MultiWoAvailable {
				updatedWC.MultiWoAvailable = backendWC.MultiWoAvailable
				needsUpdate = true
			}
			//Torn
			if cacheWC.ShiftID != backendWC.ShiftId {
				s.logger.InfoContext(ctx, "Workcenter shift configuration changed",
					slog.String("workcenter_id", backendID),
					slog.String("old_shift_id", cacheWC.ShiftID.String()),
					slog.String("new_shift_id", backendWC.ShiftId.String()),
				)
				needsUpdate = true
				updatedWC.ShiftID = backendWC.ShiftId

			}
			if needsUpdate {
				if err := s.repo.Set(ctx, backendID, updatedWC); err != nil {
					s.logger.ErrorContext(ctx, "Failed to update workcenter in cache",
						slog.String("workcenter_id", backendID),
						slog.String("error", err.Error()),
					)
					return fmt.Errorf("error updating workcenter %s: %w", backendID, err)
				}
				s.logger.InfoContext(ctx, "Workcenter updated in cache",
					slog.String("workcenter_id", backendID),
				)
			}

			
		}
	}
	s.logger.InfoContext(ctx, "BuildDTO completed %s workcenters in backend, %s in redis", len(backendWorkcenters), len(cacheWorkcenters))
	return nil
}

func(s *service) SetCurrentShift(ctx context.Context)error{
	workcenters, err := s.repo.List(ctx)
	hasChanged := false
	if err != nil {
		return fmt.Errorf("error listing workcenters: %w", err)
	}
	for _, wc := range workcenters {		
		now := time.Now()
		shiftDetail, err := s.shiftService.FindCurrentShift(ctx, now, wc.ShiftID)
		if err != nil {
			log.Printf("error finding current shift for workcenter %s: %v", wc.WorkcenterID.String(), err)
    		continue
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
			log.Default().Println("Workcenter shift updated for workcenter ", wc.WorkcenterName)
			hasChanged = true
		}
	}
	if hasChanged {
		workcenters, err := s.repo.List(ctx)
		if err != nil {
			return fmt.Errorf("error listing workcenters: %w", err)
		}
		s.hub.Broadcast("General", struct {
			Type string `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type: "Workcenter",
			Payload: workcenters,
		})
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
	duration := time.Since(start)
	
	// Record metrics
	observability.RecordShiftChange(ctx, request.WorkcenterID.String(), request.ShiftDetailId.String(), duration)
	
	log.Printf("CreateWorkcenterShift took %v", duration)
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

func (s *service) GetWorkcenterDTO(ctx context.Context, id string) (*models.WorkcenterDTO, error) {
	wc, source, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == ErrWorkcenterNotFound {
			return nil, nil
		}
		return nil, err
	}
	if source == models.SourceMemory || source == models.SourceRedis {
		return &wc, nil
	}
	return nil, nil
}

func (s *service) GetAllWorkcenters(ctx context.Context) ([]models.WorkcenterDTO, error) {
	return s.repo.List(ctx)
}