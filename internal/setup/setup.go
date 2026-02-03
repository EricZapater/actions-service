package setup

import (
	"actions-service/internal/clients"
	"actions-service/internal/config"
	"actions-service/internal/operator"
	"actions-service/internal/shift"
	"actions-service/internal/state"
	"actions-service/internal/status"
	"actions-service/internal/validator"
	"actions-service/internal/workcenter"
	"actions-service/internal/ws"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ShiftService shift.Service
	WorkcenterService workcenter.Service
	OperatorService operator.Service
	StatusService status.Service
	ValidatorService validator.Service  // ⭐ Add validator to services
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
    statusRepo := status.NewStatusRepository(state, redisClient)

	// ⭐ STEP 1: Create base services WITHOUT validator (to avoid circular dependency)
	log.Println("🔧 Creating base services...")
	operatorServiceBase := operator.NewOperatorService(client, *operatorRepo, workcenterService, hub, nil)

	// ⭐ STEP 2: Create VALIDATOR using base services (via ports)
	log.Println("✅ Creating validator service...")
	validatorService := validator.NewValidatorService(
		operatorServiceBase,  // implements OperatorPort
		statusRepo,           // implements StatusRepository
		nil,                  // no WorkOrderPhasePort (not needed yet)
		workcenterService,    // implements WorkcenterPort
	)

	// ⭐ STEP 3: Re-create services WITH validator injected
	log.Println("🔄 Re-creating services with validator...")
	operatorService := operator.NewOperatorService(client, *operatorRepo, workcenterService, hub, validatorService)
	statusService := status.NewStatusService(client, *statusRepo, workcenterService, hub, validatorService)

	log.Println("✅ All services created successfully!")

	// Create handlers
	operatorHandler := operator.NewHandler(operatorService)
	statusHandler := status.NewHandler(statusService)

	services := Services{
		ShiftService: shiftService,
		WorkcenterService: workcenterService,
		OperatorService: operatorService,
		StatusService: statusService,
		ValidatorService: validatorService,
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