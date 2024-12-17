package main

import (
	"actions-service/internal"
	"actions-service/internal/clients"
	"actions-service/internal/config"
	"actions-service/internal/controllers"
	"actions-service/internal/server"
	"actions-service/internal/services"
	"context"
	"fmt"
	"log"
	"time"
)

func Init()(*server.Services, *server.Controllers, error){
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Couldn't load environment %v", err)
		return nil, nil, err
	}
	client, err := clients.NewClientWithResponses(cfg.BackendUrl)
	if err != nil {
		log.Fatalf("Couldn't create client %v", err)
		return nil, nil, err
	}
	state := internal.GetInstance()
	statusService := services.NewStatusService(client, state)
	operatorService := services.NewOperatorService(client, state)
	services := &server.Services{
		Status: statusService,
		Operator: operatorService,
	}
	statusController := controllers.NewStatusController(statusService)
	operatorController := controllers.NewOperatorController(operatorService)
	healthCheckController := controllers.NewHealthCheckController()
	controllers := &server.Controllers{
		Status: statusController,
		Operator: operatorController,
		HealthCheck: healthCheckController,
	}
	return services, controllers, nil
}

func main() {	
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Couldn't load environment %v", err)
	}
	client, err := clients.NewClientWithResponses(cfg.BackendUrl)
	if err != nil {
		log.Fatalf("Couldn't create client %v", err)
	}
	workcenterService := services.NewWorkcenterService(client)
	startTime := time.Now()
	_ , err = workcenterService.BuildWorkcenterDTO(ctx)	
	if err != nil {
		log.Fatalf("error building DTO: %v", err)
	}
	endTime := time.Now()	
	fmt.Println(endTime.Sub(startTime))	
	
	if err != nil {
		fmt.Println("error retrieving instance %v", err)
	}	
	//state := internal.GetInstance()
	
	
	services, controllers, err := Init()
	if err != nil {
		fmt.Println("error initializing services %v", err)
	}
	server.Run(cfg, *services, *controllers)

	ticker := time.NewTicker(1*time.Minute)
	defer ticker.Stop()

	for {
		select{
		case <- ticker.C:
			fmt.Println("Ticker: ", time.Now())
		}
	}

}

