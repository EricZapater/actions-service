package operator

import (
	"actions-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ClockIn godoc
// @Summary Registra l'entrada d'un operador
// @Description Registra l'hora d'entrada d'un operador en un workcenter específic
// @Tags operator
// @Accept json
// @Produce json
// @Param request body models.OperatorRequest true "Dades del clock in"
// @Success 200 {object} models.ResponseMessage
// @Failure 400 {object} models.ResponseMessage "Request invàlid"
// @Failure 500 {object} models.ResponseMessage "Error intern del servidor"
// @Router /operator/clockin [post]
func (h *Handler) ClockIn(c *gin.Context) {
	var req models.OperatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Result:  "error",
			Message: "Invalid request",
			Content: err.Error(),
		})
		return
	}	
	if err := h.service.ClockIn(c.Request.Context(), req.OperatorID.String(), req.WorkcenterID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Result:  "error",
			Message: "Failed to register clock in",
			Content: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Clock in registered successfully",
	})
}

// ClockOut godoc
// @Summary Registra la sortida d'un operador
// @Description Registra l'hora de sortida d'un operador d'un workcenter específic
// @Tags operator
// @Accept json
// @Produce json
// @Param request body models.OperatorRequest true "Dades del clock out"
// @Success 200 {object} models.ResponseMessage
// @Failure 400 {object} models.ResponseMessage "Request invàlid"
// @Failure 500 {object} models.ResponseMessage "Error intern del servidor"
// @Router /operator/clockout [post]
func (h *Handler) ClockOut(c *gin.Context) {
	var req models.OperatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Result:  "error",
			Message: "Invalid request",
			Content: err.Error(),
		})
		return
	}
	if err := h.service.ClockOut(c.Request.Context(), req.OperatorID.String(), req.WorkcenterID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Result:  "error",
			Message: "Failed to register clock out",
			Content: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Clock out registered successfully",
	})
}