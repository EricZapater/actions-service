package operator

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	BuilDTO(ctx context.Context)error
	ClockIn(ctx context.Context, operatorID, workcenterID string)error
	ClockOut(ctx context.Context, operatorID, workcenterID string)error
}

type service struct {
	client clients.HttpBackendClient
	repo Repository	
	port WorkcenterPort
	hub *ws.Hub
}

func NewOperatorService(client clients.HttpBackendClient, repo Repository, port WorkcenterPort, hub *ws.Hub) Service {
	return &service{
		client: client,
		repo: repo,
		port: port,
		hub: hub,
	}
}

func (s *service) BuilDTO(ctx context.Context)error {
	url := "/api/Operator"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {		
		return fmt.Errorf("failed to get operators: %s", response.Status)
	}

	var operators []models.OperatorResponse
	err = json.NewDecoder(response.Body).Decode(&operators)
	if err != nil {
		return err
	}
	
	url = "/api/OperatorType"
	response, err = s.client.DoGetRequest(ctx, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return fmt.Errorf("failed to get operators: %s", response.Status)
	}

	var operatorTypes []models.OperatorTypeResponse
	err = json.NewDecoder(response.Body).Decode(&operatorTypes)
	if err != nil {
		return  err
	}	
	
	for _, operator := range operators {
		var OperatorDTO models.OperatorDTO
		OperatorDTO.OperatorID = operator.OperatorId
		OperatorDTO.OperatorCode = operator.Code
		OperatorDTO.OperatorName = operator.Name
		OperatorDTO.OperatorSurname = operator.Surname
		for _, operatorType := range operatorTypes {
			if operatorType.OperatorTypeId == operator.OperatorTypeID {
				OperatorDTO.OperatorTypeID = operatorType.OperatorTypeId
				OperatorDTO.OperatorTypeName = operatorType.Name
				OperatorDTO.OperatorTypeDescription = operatorType.Description
				OperatorDTO.OperatorTypeCost = operatorType.Cost
				break
			}
		}
		err = s.repo.Set(ctx, OperatorDTO.OperatorID.String(), OperatorDTO)
		if err != nil {
			log.Printf("error saving operator %s: %v", OperatorDTO.OperatorID.String(), err)
		}
	}

	return nil
}

func (s *service) ClockIn(ctx context.Context, operatorID, workcenterID string)error {
	wc, err := s.port.GetWorkcenterDTO(ctx, workcenterID)
	if err != nil {
		return fmt.Errorf("error checking workcenter existence: %w", err)
	}
	
	operator, _, err := s.repo.FindByID(ctx, operatorID)
	if err != nil {
		if err == ErrOperatorNotFound {
			return fmt.Errorf("operator %s not found", operatorID)
		}
		return fmt.Errorf("error finding operator %s: %w", operatorID, err)
	}
	now := time.Now()
	
	clockindto := models.OperatorClockInDTO{
		OperatorId: uuid.MustParse(operatorID),
		WorkcenterId: uuid.MustParse(workcenterID),
		Timestamp: now.Format("2006-01-02T15:04:05"),
	}
	url := "/api/WorkcenterShift/Operator/in"
	response, err := s.client.DoPostRequest(ctx, url, clockindto)
	if err != nil || response.StatusCode > 299 {
		log.Printf("Something went wrong calling the backend %v", err)
		return err
	}
	//Fer el set al repo PERO del JSON sencer del workcenter
	wc.Operators = append(wc.Operators, operator)
	if err := s.repo.SetWorkcenterDTO(ctx, wc.WorkcenterID.String(), *wc); err != nil {
		return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
	}
	s.hub.Broadcast(wc.WorkcenterID.String(), struct {
			Type string `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type: "workcenter_update",
			Payload: wc,
		})
	return nil
}

func (s *service) ClockOut(ctx context.Context, operatorID, workcenterID string)error {
	log.Printf("out")
	wc, err := s.port.GetWorkcenterDTO(ctx, workcenterID)
	if err != nil {
		return fmt.Errorf("error checking workcenter existence: %w", err)
	}
	
	_, _, err = s.repo.FindByID(ctx, operatorID)
	if err != nil {
		if err == ErrOperatorNotFound {
			return fmt.Errorf("operator %s not found", operatorID)
		}
		return fmt.Errorf("error finding operator %s: %w", operatorID, err)
	}
	now := time.Now()
	
	clockindto := models.OperatorClockInDTO{
		OperatorId: uuid.MustParse(operatorID),
		WorkcenterId: uuid.MustParse(workcenterID),
		Timestamp: now.Format("2006-01-02T15:04:05"),
	}
	url := "/api/WorkcenterShift/Operator/out"
	response, err := s.client.DoPostRequest(ctx, url, clockindto)
	if err != nil || response.StatusCode > 299 {
		log.Printf("Something went wrong calling the backend %v", err)
		return err
	}
	operators := wc.Operators
	filtered := make([]models.OperatorDTO, 0, len(operators))
	for _, op := range operators {
		if op.OperatorID.String() != operatorID {
			filtered = append(filtered, op)
		}
	}
	wc.Operators = filtered

	dump, _ := json.MarshalIndent(wc, "", "  ")
	log.Printf("Workcenter after clockout: %s", string(dump))

	if err := s.repo.SetWorkcenterDTO(ctx, wc.WorkcenterID.String(), *wc); err != nil {
		return fmt.Errorf("error updating workcenter %s: %w", wc.WorkcenterID.String(), err)
	}
	s.hub.Broadcast(wc.WorkcenterID.String(), struct {
			Type string `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type: "workcenter_update",
			Payload: wc,
		})
	return nil
}