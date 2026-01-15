package setup

import (
	"actions-service/internal/clients"
	"actions-service/internal/config"
	"actions-service/internal/operator"
	"actions-service/internal/shift"
	"actions-service/internal/status"
	"actions-service/internal/workcenter"
	"actions-service/internal/workorderphase"
	"actions-service/internal/ws"
	"context"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ShiftService shift.Service
	WorkcenterService workcenter.Service
	OperatorService operator.Service
	StatusService status.Service
	WorkOrderPhaseService workorderphase.Service
}

type Handlers struct {
	OperatorHandler *operator.Handler
	StatusHandler *status.Handler
	WorkOrderPhaseHandler *workorderphase.Handler
}

type App struct {
	Cfg *config.Config	
	Services Services
	Handlers Handlers
	Hub *ws.Hub
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	client := clients.NewHttpBackendClient(cfg.BackendUrl)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: "", 
		DB:       0,  
	})

	hub := ws.NewHub()	
	
	shiftRepo := shift.NewShiftRepository()
	shiftService := shift.NewShiftService(client, *shiftRepo)

	workcenterRepo := workcenter.NewWorkcenterRepository(redisClient)
	workcenterService := workcenter.NewWorkcenterService(client, *workcenterRepo, shiftService, hub)

	operatorRepo := operator.NewOperatorRepository(redisClient)
	operatorService := operator.NewOperatorService(client, *operatorRepo, workcenterService, hub)
	operatorHandler := operator.NewHandler(operatorService)

    statusRepo := status.NewStatusRepository(redisClient)
    statusService := status.NewStatusService(client, *statusRepo, workcenterService, operatorService, hub)
	statusHandler := status.NewHandler(statusService)

	workorderphaseRepo := workorderphase.NewWorkOrderPhaseRepository(redisClient)
	workorderphaseService := workorderphase.NewWorkOrderPhaseService(client, *workorderphaseRepo, workcenterService, hub, statusService, operatorService)
	workorderphaseHandler := workorderphase.NewHandler(workorderphaseService)

	services := Services{
		ShiftService: shiftService,
		WorkcenterService: workcenterService,
		OperatorService: operatorService,
		StatusService: statusService,
		WorkOrderPhaseService: workorderphaseService,
	}

	handlers := Handlers{
		OperatorHandler: operatorHandler,
		StatusHandler: statusHandler,
		WorkOrderPhaseHandler: workorderphaseHandler,
	}

	return &App{
		Cfg: cfg,		
		Services: services,
		Handlers: handlers,
		Hub: hub,
	}, nil
}