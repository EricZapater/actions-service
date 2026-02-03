package bootstrap

import (
	"actions-service/internal/clients"
	"context"
	"fmt"
)

type Service struct {
	RedisRepo *RedisRepo
	client    clients.HttpBackendClient	
}

type service interface {
	InitDTO(ctx context.Context) error
}

func NewService(redisRepo *RedisRepo, client clients.HttpBackendClient) *Service {
	return &Service{
		RedisRepo: redisRepo,
		client:    client,
	}
}

func (s *Service) InitDTO(ctx context.Context) error {
	url := "/api/WorkcenterShift/Currents"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {		
		return fmt.Errorf("failed to get operators: %s", response.Status)
	}
	return nil
}
