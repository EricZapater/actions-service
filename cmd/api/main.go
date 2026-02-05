package main

import (
	"actions-service/internal/observability"
	"actions-service/internal/server"
	"actions-service/internal/setup"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main() {	
	startTime := time.Now()
	ctx := context.Background()
	app, err := setup.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize observability
	shutdown, err := observability.InitTelemetry(ctx, observability.Config{
		OtelEndpoint:   app.Cfg.OtelEndpoint,
		ServiceName:    app.Cfg.ServiceName,
		ServiceVersion: app.Cfg.ServiceVersion,
		Environment:    app.Cfg.Environment,
	})
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()

	// Initialize metrics
	if err := observability.InitMetrics(); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	app.Services.ShiftService.BuildDTO(ctx)
	app.Services.OperatorService.BuilDTO(ctx)
	app.Services.StatusService.BuildDTO(ctx)
	app.Services.BootstrapService.InitDTO(ctx)	
	app.Services.WorkcenterService.BuildDTO(ctx)	
	
	
	
	
	go server.Run(app)
	
	endTime := time.Now()	
	log.Printf("Startup time: %v", endTime.Sub(startTime))

	ticker := time.NewTicker(1*time.Minute)
	defer ticker.Stop()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				//workcenterService.CheckWorkcenterShift(ctx)
				app.Services.WorkcenterService.SetCurrentShift(ctx)
			case <-stopChan:
				log.Println("Shutting down ticker...")
				return
			}
		}
	}()
	
	<-stopChan // Espera senyal per aturar l'aplicació
	log.Println("Application stopped")
}

