package server

import (
	"actions-service/internal/models"
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

// HealthCheck godoc
// @Summary Comprova l'estat del servei
// @Description Retorna l'estat de salut del servidor i l'estat actual del sistema
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMessage
// @Router /healthcheck [get]
func (h *Handler) HealthCheck(ctx *gin.Context) {
	s := h.app.State.GetState()
	ctx.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Service is healthy",
		Content: s,
	})
}

// ReloadDTO godoc
// @Summary Recarrega els DTOs del sistema
// @Description Recarrega els DTOs de shifts, operadors i workcenters i retorna l'estat actualitzat
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMessage
// @Router /reload [get]
func (h *Handler) ReloadDTO(ctx *gin.Context) {
	h.app.Services.ShiftService.BuildDTO(ctx)
	h.app.Services.OperatorService.BuilDTO(ctx)
	h.app.Services.WorkcenterService.BuildDTO(ctx)	
	s := h.app.State.GetState()
	ctx.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "DTOs reloaded successfully",
		Content: s,
	})
}

// WSGeneral godoc
// @Summary WebSocket de connexió general
// @Description Estableix una connexió WebSocket per rebre actualitzacions generals del sistema
// @Tags websocket
// @Accept json
// @Produce json
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/general [get]
func (h *Handler) WSGeneral(ctx *gin.Context) {
	state := h.app.State.GetState()
	h.app.Hub.HandleWS(ctx.Writer, ctx.Request, "general", state)
}

// WSWorkcenter godoc
// @Summary WebSocket de connexió per workcenter
// @Description Estableix una connexió WebSocket per rebre actualitzacions específiques d'un workcenter
// @Tags websocket
// @Accept json
// @Produce json
// @Param id path string true "ID del Workcenter"
// @Success 101 {string} string "Switching Protocols"
// @Failure 404 {object} models.ResponseMessage "Workcenter no trobat"
// @Router /ws/workcenter/{id} [get]
func(h *Handler) WSWorkcenter(ctx *gin.Context) {
	workcenterID := ctx.Param("id")
	state := h.app.State.GetState()	
	h.app.Hub.HandleWS(ctx.Writer, ctx.Request, workcenterID, state.Workcenters[workcenterID])
}