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
	client clients.HttpBackendClient
}

func NewAreaService(client clients.HttpBackendClient) AreaService{
	return &areaService{client: client}
}

func(s *areaService) GetAreas(ctx context.Context)([]models.Area, error){		
	response, err := s.client.DoGetRequest(ctx, "GET", "/api/area")
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
		var areas []models.Area
		if err := json.NewDecoder(response.Body).Decode(&areas); err != nil {
			return nil, fmt.Errorf("error deserializing response: %w", err)
		}
		return areas, nil
	}else{
		return nil, fmt.Errorf("Response error: %v", response.Status)
	}
	
}