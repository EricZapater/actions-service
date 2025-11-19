package setup

import (
	"actions-service/internal/clients"
	"actions-service/internal/config"
	"actions-service/internal/operator"
	"actions-service/internal/shift"
	"actions-service/internal/state"
	"actions-service/internal/status"
	"actions-service/internal/workcenter"
	"actions-service/internal/ws"
	"context"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ShiftService shift.Service
	WorkcenterService workcenter.Service
	OperatorService operator.Service
	StatusService status.Service
}

type Handlers struct {
	OperatorHandler *operator.Handler
	StatusHandler *status.Handler
}

type App struct {
	Cfg *config.Config	
	Services Services
	Handlers Handlers
	State *state.State
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
	state := state.New()

	hub := ws.NewHub()	
	

	shiftRepo := shift.NewShiftRepository(state)
	shiftService := shift.NewShiftService(client, *shiftRepo)

	workcenterRepo := workcenter.NewWorkcenterRepository(state, redisClient)
	workcenterService := workcenter.NewWorkcenterService(client, *workcenterRepo, shiftService, hub)

	operatorRepo := operator.NewOperatorRepository(state, redisClient)
	operatorService := operator.NewOperatorService(client, *operatorRepo, workcenterService, hub)
	operatorHandler := operator.NewHandler(operatorService)

    statusRepo := status.NewStatusRepository(state, redisClient)
    statusService := status.NewStatusService(client, *statusRepo, workcenterService, hub)
	statusHandler := status.NewHandler(statusService)

	services := Services{
		ShiftService: shiftService,
		WorkcenterService: workcenterService,
		OperatorService: operatorService,
		StatusService: statusService,
	}

	handlers := Handlers{
		OperatorHandler: operatorHandler,
		StatusHandler: statusHandler,
	}

	return &App{
		Cfg: cfg,		
		Services: services,
		Handlers: handlers,
		State: state,
		Hub: hub,
	}, nil
}