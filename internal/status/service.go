package status

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
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
}

type service struct {
	client clients.HttpBackendClient	
    hub *ws.Hub
    repo Repository
    workcenterPort WorkcenterPort
    operatorPort OperatorPort
}

func NewStatusService(client clients.HttpBackendClient, repo Repository, workcenterPort WorkcenterPort, operatorPort OperatorPort, hub *ws.Hub) Service{
	return &service{
		client: client,
        repo: repo,
        workcenterPort: workcenterPort,
        operatorPort: operatorPort,
		hub: hub,
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

	url = "/api/WorkcenterCost"
	responseCost, err := s.client.DoGetRequest(ctx, url)	
	if err != nil {
		return err
	}
	defer responseCost.Body.Close()
	if responseCost.StatusCode > 299 {
		return fmt.Errorf("failed to get status costs: %s", responseCost.Status)
	}

	var statuscosts []models.StatusCostResponse
	err = json.NewDecoder(responseCost.Body).Decode(&statuscosts)
	if err != nil {
		return err
	}
    for _, cost := range statuscosts {
        var dto models.StatusDTO
        dto.WorkcenterId = cost.WorkcenterId
        dto.StatusId = cost.StatusId
        dto.Cost = cost.Cost
        for _, st := range statuses {
            if st.StatusId == cost.StatusId {
                dto.Description = st.Description
                dto.Closed = st.Closed
                dto.Color = st.Color
                dto.OperatorsAllowed = st.OperatorsAllowed
                dto.Stopped = st.Stopped
                break
            }
        }
        // composite key to avoid collisions
        key := fmt.Sprintf("%s:%s", dto.WorkcenterId.String(), dto.StatusId.String())
        if err := s.repo.Set(ctx, key, dto); err != nil {
            return err
        }
    }

    return nil
}

func (s *service) StatusIn(ctx context.Context, workcenterID, statusID string, reasonID *string) error {
    wc, err := s.workcenterPort.GetWorkcenterDTO(ctx, workcenterID)
    if err != nil {
        return fmt.Errorf("error checking workcenter existence: %w", err)
    }
    if wc == nil {
        return fmt.Errorf("workcenter %s not found", workcenterID)
    }

    key := fmt.Sprintf("%s:%s", workcenterID, statusID)
    st, _, err := s.repo.FindByID(ctx, key)
    if err != nil {
        return fmt.Errorf("status %s for workcenter %s not found: %w", statusID, workcenterID, err)
    }

    if !st.OperatorsAllowed {
        //operators out
        for _, operator := range wc.Operators {            
            s.operatorPort.ClockOut(ctx, operator.OperatorID.String(), workcenterID)
        }
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
    wc.StatusReasonId = st.StatusReasonId
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

    return nil
}