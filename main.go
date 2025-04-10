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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init()(*server.Services, *server.Controllers, error){
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Couldn't load environment %v", err)
		return nil, nil, err
	}
	client := clients.NewHttpBackendClient(cfg.BackendUrl)
	
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
	client := clients.NewHttpBackendClient(cfg.BackendUrl)
	
	shiftService := services.NewShiftService(client)
	_, err = shiftService.BuildShiftsDTO(ctx)
	if err != nil {
		log.Fatalf("error building shift DTO: %v", err)
	}
	workcenterService := services.NewWorkcenterService(client)
	startTime := time.Now()
	_ , err = workcenterService.BuildWorkcenterDTO(ctx)	
	if err != nil {
		log.Fatalf("error building workcenter DTO: %v", err)
	}
	endTime := time.Now()	
	fmt.Println("starttime:",  endTime.Sub(startTime))	
	
	if err != nil {
		log.Fatalf("error retrieving instance %v", err)
	}	
	
	
	
	services, controllers, err := Init()
	if err != nil {
		log.Fatalf("error initializing services %v", err)
	}
	go server.Run(cfg, *services, *controllers)


	ticker := time.NewTicker(1*time.Minute)
	defer ticker.Stop()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				//fmt.Println("ticker")
				workcenterService.CheckWorkcenterShift(ctx)
			case <-stopChan:
				log.Println("Shutting down ticker...")
				return
			}
		}
	}()
	
	<-stopChan // Espera senyal per aturar l'aplicació
	log.Println("Application stopped")
	

}

