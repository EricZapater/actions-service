package controllers

import (
	"actions-service/internal"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HealthCheckController struct{}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) HealthCheck(ctx *gin.Context) {
	state := internal.GetInstance()
	state.Mu.Lock()
	defer state.Mu.Unlock()
	wcid, err := uuid.Parse("aa86ece8-4a67-4099-adaa-16f543d75f9d") 
	if err != nil {
		message:=fmt.Sprintf("error parsing guid, error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}
	fmt.Println(state.Workcenters[wcid])
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Server Running",
		"state" : state.Workcenters,
		"shifts": state.Shifts,
	})
}