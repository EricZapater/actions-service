package services

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type AreaService interface {
	GetAreas(ctx context.Context)([]models.Area, error)
}

type areaService struct {
	client *clients.ClientWithResponses
}

func NewAreaService(client *clients.ClientWithResponses) AreaService{
	return &areaService{client: client}
}

func(s *areaService) GetAreas(ctx context.Context)([]models.Area, error){		
	response, err := s.client.GetApiAreaWithResponse(ctx)
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
		var areas []models.Area
		if err := json.Unmarshal(response.Body, &areas); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return areas, nil
	}else{
		return nil, fmt.Errorf("Response error: %v", response.HTTPResponse.Status)
	}
	
}