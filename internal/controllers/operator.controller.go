package controllers

import (
	"actions-service/internal"
	"actions-service/internal/models"
	"actions-service/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type OperatorController struct {
	service services.OperatorService
}

func NewOperatorController(service services.OperatorService) *OperatorController{
	return &OperatorController{service: service}
}

func(c *OperatorController)ClockIn(ctx *fiber.Ctx)error{
	var request models.OperatorRequest
	if err := ctx.BodyParser(&request); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
	fmt.Println(state)
	fmt.Println("clock in: ",request)
	return ctx.Status(fiber.StatusAccepted).JSON(fmt.Sprintf("%v,\n%v", request, state))
}

func(c *OperatorController)ClockOut(ctx *fiber.Ctx)error{
	var request models.OperatorRequest
	if err := ctx.BodyParser(&request); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }
	fmt.Println("clock in: ",request)
	return ctx.Status(fiber.StatusAccepted).JSON(request)
}