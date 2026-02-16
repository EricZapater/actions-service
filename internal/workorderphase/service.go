package workorderphase

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	WorkOrderPhaseIn(ctx context.Context, req models.WorkOrderPhaseAndStatusRequest)error
	WorkOrderPhaseOut(ctx context.Context, req models.WorkOrderPhaseOutRequest)error
	FindByID(ctx context.Context, id string) (models.WorkOrderPhaseResponse, error)
}

type service struct {
	client clients.HttpBackendClient
	hub *ws.Hub	
	repo Repository
	workcenterPort WorkcenterPort	
	statusPort StatusPort
	operatorPort OperatorPort
}

func NewWorkOrderPhaseService(client clients.HttpBackendClient, repo Repository, workcenterPort WorkcenterPort, hub *ws.Hub, statusPort StatusPort, operatorPort OperatorPort) Service{
	return &service{
		client: client,
		hub: hub,
		repo: repo,
		workcenterPort: workcenterPort,
		statusPort: statusPort,
		operatorPort: operatorPort,
	}
}

func (s *service) WorkOrderPhaseIn(ctx context.Context, req models.WorkOrderPhaseAndStatusRequest)error{
	wc, err := s.workcenterPort.GetWorkcenterDTO(ctx, req.WorkcenterID)
    if err != nil {
        return fmt.Errorf("error checking workcenter existence: %w", err)
    }
    if wc == nil {
        return NewServiceError(http.StatusNotFound, fmt.Sprintf("workcenter %s not found", req.WorkcenterID), nil)
    }

	if !wc.MultiWoAvailable && len(wc.WorkOrders) > 0 {
		return NewServiceError(http.StatusConflict, fmt.Sprintf("workcenter %s already has a workorder", req.WorkcenterID), nil)
	}

	//Get Phase information
	phase, err := s.FindByID(ctx, req.WorkOrderPhaseId)
	if err != nil {
		return err
	}
	now := time.Now().Format("2006-01-02T15:04:05")
	wo := s.createWorkOrderDTO(req.WorkOrderPhaseId, phase, now)
	
	//Comprovar si en la request hi ha MachineStatusId
	request := models.WorkOrderPhaseAndStatusRequest{}
	st := models.StatusDTO{}
	if req.MachineStatusId == nil{
		request.WorkcenterID = req.WorkcenterID
		request.WorkOrderPhaseId = req.WorkOrderPhaseId		
		request.TimeStamp = &now
		response, err := s.client.DoPostRequest(ctx, "/api/WorkcenterShift/WorkOrderPhase/In", request)
		if err != nil || response == nil || response.StatusCode > 299 {
			var status string
			var code int
			if response != nil {
				status = response.Status
				code = response.StatusCode
			}
			return fmt.Errorf("backend workorderphase in failed (code %d, status %s): %w", code, status, err)
		}
		//fmt.Println(response)
		if response != nil && response.Body != nil { _ = response.Body.Close() }
	} else {
		st, err = s.statusPort.FindByID(ctx, *req.MachineStatusId)
		if err != nil {
			return fmt.Errorf("error checking status existence: %w", err)
		}
		if !st.OperatorsAllowed {
        //operators out
			for _, operator := range wc.Operators {            
				s.operatorPort.ClockOut(ctx, operator.OperatorID.String(), req.WorkcenterID)
			}
			wc.Operators = []models.OperatorDTO{}
		}	
        
		request.WorkcenterID = req.WorkcenterID
		request.WorkOrderPhaseId = req.WorkOrderPhaseId		
		request.MachineStatusId = req.MachineStatusId	
		request.TimeStamp = &now
		response, err := s.client.DoPostRequest(ctx, "/api/WorkcenterShift/WorkOrderPhaseAndStatus/In", request)
		if err != nil || response == nil || response.StatusCode > 299 {
			var status string
			var code int
			if response != nil {
				status = response.Status
				code = response.StatusCode
			}
			return fmt.Errorf("backend workorderphase and status in failed (code %d, status %s): %w", code, status, err)
		}

		if response != nil && response.Body != nil { _ = response.Body.Close() }
    	
		wc.StatusID = st.StatusId
    	wc.StatusReasonId = &uuid.Nil
    	wc.StatusName = st.Description
    	wc.StatusOperatorsAllowed = st.OperatorsAllowed
    	wc.StatusClosed = st.Closed
    	wc.StatusStopped = st.Stopped
    	wc.StatusColor = st.Color
    	wc.StatusStartTime = time.Now()
	}
	
	exists := false
	for _, workorder := range wc.WorkOrders {
		if workorder.WorkOrderPhaseId == req.WorkOrderPhaseId {
			exists = true
			break
		}
	}

	if !exists {
		wc.WorkOrders = append(wc.WorkOrders, wo)
	}

	if err := s.repo.SetWorkcenterDTO(ctx, wc.WorkcenterID.String(), *wc); err != nil {
        return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
    }
	
	s.publishWorkcenter(ctx, *wc)
		
	return nil	
}

func (s *service) WorkOrderPhaseOut(ctx context.Context, req models.WorkOrderPhaseOutRequest)error{
 	wc, err := s.workcenterPort.GetWorkcenterDTO(ctx, req.WorkcenterID)
    if err != nil {
        return fmt.Errorf("error checking workcenter existence: %w", err)
    }
    if wc == nil {
        return fmt.Errorf("workcenter %s not found", req.WorkcenterID)
    }
	//comprovar quantitats
	if req.QuantityOk != nil && req.QuantityKo != nil {
		quantityRequest := models.WorkOrderPhaseQuantitiesRequest{
			WorkcenterID:     req.WorkcenterID,
			WorkOrderPhaseId: req.WorkOrderPhaseId,
			QuantityOk:       *req.QuantityOk,
			QuantityKo:       *req.QuantityKo,
		}
		err = s.updateQuantities(ctx, quantityRequest)
		if err != nil {
			return err
		}
	}
	// Prepare backend request
	backendRequest := models.BackendWorkOrderPhaseOutRequest{
		WorkcenterID:      req.WorkcenterID,
		WorkOrderPhaseId:  req.WorkOrderPhaseId,
		TimeStamp:         time.Now().Format("2006-01-02T15:04:05"),
		WorkOrderStatusId: req.WorkOrderStatusId,
		NextWorkOrderPhaseId: req.NextWorkOrderPhaseId,
		NextMachineStatusId: req.NextMachineStatusId,
	}
	if req.NextMachineStatusId != nil {
		st, err := s.statusPort.FindByID(ctx, *req.NextMachineStatusId)
		if err != nil {
			return err
		}
		if !st.OperatorsAllowed {
			for _, operator := range wc.Operators {
				s.operatorPort.ClockOut(ctx, operator.OperatorID.String(), req.WorkcenterID)
			}
			wc.Operators = []models.OperatorDTO{}
		}
		wc.StatusID = st.StatusId
		wc.StatusReasonId = &uuid.Nil
		wc.StatusName = st.Description
		wc.StatusOperatorsAllowed = st.OperatorsAllowed
		wc.StatusClosed = st.Closed
		wc.StatusStopped = st.Stopped
		wc.StatusColor = st.Color
		wc.StatusStartTime = time.Now()
	}
		
	response, err := s.client.DoPostRequest(ctx, "/api/WorkcenterShift/WorkOrderPhase/Out", backendRequest)
		if err != nil || response == nil || response.StatusCode > 299 {
			var status string
			var code int
			if response != nil {
				status = response.Status
				code = response.StatusCode
			}
			return fmt.Errorf("backend workorderphase in failed (code %d, status %s): %w", code, status, err)
		}
		
	if response != nil && response.Body != nil { _ = response.Body.Close() }
	workorders := wc.WorkOrders
	filtered := make([]models.WorkOrderDTO, 0, len(workorders))
	for _, workorder := range workorders {
		if workorder.WorkOrderPhaseId != req.WorkOrderPhaseId {
			filtered = append(filtered, workorder)
		}
	}
	wc.WorkOrders = filtered
	if req.NextWorkOrderPhaseId != nil {
		phase, err := s.FindByID(ctx, *req.NextWorkOrderPhaseId)
		if err != nil {
			return err
		}
		
		now := time.Now().Format("2006-01-02T15:04:05")
		wo := s.createWorkOrderDTO(*req.NextWorkOrderPhaseId, phase, now)
		wc.WorkOrders = append(wc.WorkOrders, wo)
	}
	
	if err := s.repo.SetWorkcenterDTO(ctx, wc.WorkcenterID.String(), *wc); err != nil {
        return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
    }
	
	s.publishWorkcenter(ctx, *wc)
		
	return nil	
}

func(s *service) updateQuantities(ctx context.Context, req models.WorkOrderPhaseQuantitiesRequest) error {
	response, err := s.client.DoPutRequest(ctx, "/api/WorkcenterShift/WorkOrderPhase/Quantities", req)
		if err != nil || response == nil || response.StatusCode > 299 {
			var status string
			var code int
			if response != nil {
				status = response.Status
				code = response.StatusCode
			}
			return fmt.Errorf("backend workorderphase in failed (code %d, status %s): %w", code, status, err)
		}		
		if response != nil && response.Body != nil { _ = response.Body.Close() }
		return nil
}


func(s *service) publishWorkcenter(ctx context.Context, wc models.WorkcenterDTO) error {
		s.hub.Broadcast(wc.WorkcenterID.String(), struct {
				Type string `json:"type"`
				Payload interface{} `json:"payload"`
			}{
				Type: "Workcenter",
				Payload: wc,
			})
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
			return nil
}

func(s *service) FindByID(ctx context.Context, id string) (models.WorkOrderPhaseResponse, error) {
	response, err := s.client.DoGetRequest(ctx, "/api/WorkOrder/Detailed/ByWorkOrderPhase/" + id)
	if err != nil || response == nil || response.StatusCode > 299 {
		var status string
		var code int
		if response != nil {
			status = response.Status
			code = response.StatusCode
		}
		return models.WorkOrderPhaseResponse{}, fmt.Errorf("backend workorderphase in failed (code %d, status %s): %w", code, status, err)
	}
	var phase []models.WorkOrderPhaseResponse
	err = json.NewDecoder(response.Body).Decode(&phase)
	if err != nil {
		return models.WorkOrderPhaseResponse{}, err
	}
	if response != nil && response.Body != nil { _ = response.Body.Close() }
	return phase[0], nil
}


func (s *service) createWorkOrderDTO(phaseid string, phase models.WorkOrderPhaseResponse, now string) models.WorkOrderDTO {
    return models.WorkOrderDTO{
        WorkOrderPhaseId: phaseid,
        StartTime: now,
        PlannedQuantity: phase.PlannedQuantity,
        WorkOrderCode: phase.WorkOrderCode,
        WorkOrderPhaseCode: phase.WorkOrderPhaseCode,
        WorkOrderPhaseDescription: phase.WorkOrderPhaseDescription,
        ReferenceCode: phase.ReferenceCode,
        ReferenceDescription: phase.ReferenceDescription,
    }
}