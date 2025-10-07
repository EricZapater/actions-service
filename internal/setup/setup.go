package setup

import (
	"actions-service/internal/clients"
	"actions-service/internal/config"
	"actions-service/internal/shift"
	"actions-service/internal/state"
	"actions-service/internal/workcenter"
	"context"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ShiftService shift.Service
	WorkcenterService workcenter.Service
}

type App struct {
	Cfg *config.Config	
	Services Services
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

	shiftRepo := shift.NewShiftRepository(state)
	shiftService := shift.NewShiftService(client, *shiftRepo)

	workcenterRepo := workcenter.NewWorkcenterRepository(state, redisClient)
	workcenterService := workcenter.NewWorkcenterService(client, *workcenterRepo, shiftService)

	services := Services{
		ShiftService: shiftService,
		WorkcenterService: workcenterService,
	}
	return &App{
		Cfg: cfg,		
		Services: services,
	}, nil
}