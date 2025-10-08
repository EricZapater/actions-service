package main

import (
	"actions-service/internal/server"
	"actions-service/internal/setup"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main() {	
	ctx := context.Background()
	app, err := setup.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	app.Services.ShiftService.BuildDTO(ctx)
	app.Services.OperatorService.BuilDTO(ctx)
	app.Services.WorkcenterService.BuildDTO(ctx)	
	
	
	go server.Run(app)


	ticker := time.NewTicker(1*time.Minute)
	defer ticker.Stop()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("ticker")
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

