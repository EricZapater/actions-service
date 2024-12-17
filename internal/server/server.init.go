package server

import (
	"actions-service/internal/controllers"
	"actions-service/internal/services"
)

type Controllers struct {
	Status *controllers.StatusController
	Operator *controllers.OperatorController
	HealthCheck *controllers.HealthCheckController
}
type Services struct {
	Status services.StatusService
	Operator services.OperatorService
}