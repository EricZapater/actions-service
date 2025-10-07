package server

import (
	"actions-service/internal/setup"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *setup.App
}

func NewHandler(app *setup.App) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) ReloadDTO(ctx *gin.Context) {
	h.app.Services.ShiftService.BuildDTO(ctx)
	h.app.Services.OperatorService.BuilDTO(ctx)
	h.app.Services.WorkcenterService.BuildDTO(ctx)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}