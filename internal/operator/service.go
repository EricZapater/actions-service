package operator

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type Service interface {
	BuilDTO(ctx context.Context)error
	ClockIn(ctx context.Context, operatorID, workcenterID string)error
	ClockOut(ctx context.Context, operatorID, workcenterID string)error
}

type service struct {
	client clients.HttpBackendClient
	repo Repository	
}

func NewOperatorService(client clients.HttpBackendClient, repo Repository) Service {
	return &service{
		client: client,
		repo: repo,
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
	return nil
}

func (s *service) ClockOut(ctx context.Context, operatorID, workcenterID string)error {
	return nil
}