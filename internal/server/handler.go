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
	s := h.app.State.GetState()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"state": s,
	})
}

func (h *Handler) ReloadDTO(ctx *gin.Context) {
	h.app.Services.ShiftService.BuildDTO(ctx)
	h.app.Services.OperatorService.BuilDTO(ctx)
	h.app.Services.WorkcenterService.BuildDTO(ctx)	
	s := h.app.State.GetState()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"state": s,
	})
}

func (h *Handler) WSGeneral(ctx *gin.Context) {
	state := h.app.State.GetState()
	h.app.Hub.HandleWS(ctx.Writer, ctx.Request, "general", state)
}

func(h *Handler) WSWorkcenter(ctx *gin.Context) {
	workcenterID := ctx.Param("id")
	state := h.app.State.GetState()	
	h.app.Hub.HandleWS(ctx.Writer, ctx.Request, workcenterID, state.Workcenters[workcenterID])
}