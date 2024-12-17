package controllers

import (
	"actions-service/internal"
	"actions-service/internal/models"
	"actions-service/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusController struct {
	service services.StatusService
}

func NewStatusController(service services.StatusService) *StatusController{
	return &StatusController{service: service}
}

func (c *StatusController)UpdateWorkcenterStatus(ctx *gin.Context) {
	var status models.ChangeStatusRequest
	if err := ctx.ShouldBindJSON(&status); err != nil {
		// Si hi ha un error, retornem un error amb codi 400 (Bad Request)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	if err := c.service.UpdateWorkcenterStatus(ctx, status.WorkcenterId, status.StatusId); err != nil {
		message:=fmt.Sprintf("could not change status, error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
    ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Server Running",
		"state" : state.Workcenters,
	})
	return
}
