package status

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/validator"
	"actions-service/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	BuildDTO(ctx context.Context)error
    StatusIn(ctx context.Context, workcenterID, statusID string, reasonID *string) error
    FindByID(ctx context.Context, statusID string) (models.StatusDTO, error)
    GetDefaultStatus(ctx context.Context) (models.StatusDTO, error)
}

type service struct {
	client    clients.HttpBackendClient	
    hub       *ws.Hub
    repo      Repository
    port      WorkcenterPort
    validator validator.Service  // ⭐ Validator instead of direct operator dependency
}

func NewStatusService(client clients.HttpBackendClient, repo Repository, port WorkcenterPort, hub *ws.Hub, validator validator.Service) Service{
	return &service{
		client:    client,
        repo:      repo,
        port:      port,
		hub:       hub,
		validator: validator,
	}
}

func(s *service) BuildDTO(ctx context.Context)error{
	url := "/api/MachineStatus"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
    if response.StatusCode > 299 {
        return fmt.Errorf("failed to get statuses: %s", response.Status)
    }

    var statuses []models.StatusResponse
	err = json.NewDecoder(response.Body).Decode(&statuses)
	if err != nil {
		return err
	}

	err = s.repo.DeleteAll(ctx)
    if err != nil {
        return fmt.Errorf("error deleting statuses: %v", err)
    }
    for _, status := range statuses {
        var dto models.StatusDTO
        dto.StatusId = status.StatusId
        dto.Description = status.Description
        dto.Closed = status.Closed        
        dto.Color = status.Color
        dto.OperatorsAllowed = status.OperatorsAllowed
        dto.IsDefault = status.IsDefault
        dto.Stopped = status.Stopped
        
        
        key := fmt.Sprintf("%s", dto.StatusId.String())
        if err := s.repo.Set(ctx, key, dto); err != nil {
            return err
        }
        if dto.IsDefault {
            fmt.Println("Default status: ", dto.Description)
            if err := s.repo.SetDefault(ctx, key, dto); err != nil {
                return err
            }
        }
    }
    return nil
}

func (s *service) StatusIn(ctx context.Context, workcenterID, statusID string) error {
    // ⭐ VALIDATE AND EXECUTE: Let validator handle business rules
    // This will automatically remove operators/workorders if new status doesn't allow them
    if err := s.validator.ValidateStatusChange(ctx, workcenterID, statusID); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Get workcenter (already updated by validator if needed)
    wc, err := s.port.GetWorkcenterDTO(ctx, workcenterID)
    if err != nil {
        return fmt.Errorf("error checking workcenter existence: %w", err)
    }
    if wc == nil {
        return fmt.Errorf("workcenter %s not found", workcenterID)
    }

    //key := fmt.Sprintf("%s:%s", workcenterID, statusID)
    //st, _, err := s.repo.FindByID(ctx, key)
    st, err := s.FindByID(ctx, statusID)
    if err != nil {
        return fmt.Errorf("status %s not found: %w", statusID, err)
    }


    if !st.OperatorsAllowed {
        //operators out
        for _, operator := range wc.Operators {            
            s.operatorPort.ClockOut(ctx, operator.OperatorID.String(), workcenterID)
        }
        // Clear operators from memory to avoid overwriting the ClockOut changes
        wc.Operators = []models.OperatorDTO{}
    }

    // backend call
    req := models.StatusInRequest{}
    req.WorkcenterID = wc.WorkcenterID
    parsed, err := uuid.Parse(statusID)
    if err != nil { return fmt.Errorf("invalid statusID %s: %w", statusID, err) }
    req.StatusID = parsed
    if reasonID != nil {
        parsedReason, err := uuid.Parse(*reasonID)
        if err != nil { return fmt.Errorf("invalid reasonID %s: %w", *reasonID, err) }
        req.StatusReasonId = &parsedReason
    }       
    req.Timestamp = time.Now().Format("2006-01-02T15:04:05")
    response, err := s.client.DoPostRequest(ctx, "/api/WorkcenterShift/Workcenter/ChangeStatus", req)
    if err != nil || response == nil || response.StatusCode > 299 {
        var status string
        var code int
        if response != nil {
            status = response.Status
            code = response.StatusCode
        }
        return fmt.Errorf("backend status in failed (code %d, status %s): %w", code, status, err)
    }
    if response != nil && response.Body != nil { _ = response.Body.Close() }

    // update workcenter status fields
    wc.StatusID = st.StatusId
    wc.StatusReasonId = req.StatusReasonId
    wc.StatusName = st.Description
    wc.StatusOperatorsAllowed = st.OperatorsAllowed
    wc.StatusClosed = st.Closed
    wc.StatusStopped = st.Stopped
    wc.StatusColor = st.Color
    wc.StatusStartTime = time.Now()

    if err := s.repo.SetWorkcenterDTO(ctx, wc.WorkcenterID.String(), *wc); err != nil {
        return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
    }

    // broadcast updates
    s.hub.Broadcast(wc.WorkcenterID.String(), struct {
        Type string `json:"type"`
        Payload interface{} `json:"payload"`
    }{
        Type: "Workcenter",
        Payload: wc,
    })

	workcenters, err := s.workcenterPort.GetAllWorkcenters(ctx)
	if err != nil {
		return fmt.Errorf("error listing workcenters: %w", err)
	}

	s.hub.Broadcast("general", struct {
			Type string `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type: "Workcenter",
			Payload: workcenters,
		})

    return nil
}

func (s *service) FindByID(ctx context.Context, statusID string) (models.StatusDTO, error){
	key := fmt.Sprintf("%s", statusID)
    st, _, err := s.repo.FindByID(ctx, key)
    if err != nil {
        return models.StatusDTO{}, fmt.Errorf("status %s not found: %w", statusID, err)
    }
    return st, nil
}

func(s *service) GetDefaultStatus(ctx context.Context) (models.StatusDTO, error){
	st, err := s.repo.GetDefaultStatus(ctx)
    if err != nil {
        return models.StatusDTO{}, fmt.Errorf("default status not found: %w", err)
    }
    return st, nil
}
