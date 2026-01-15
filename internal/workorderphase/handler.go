package workorderphase

import (
	"actions-service/internal/models"
	"actions-service/internal/observability"
	"errors"
	"fmt"
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

func (h *Handler) WorkOrderPhaseIn(c *gin.Context) {
	var req models.WorkOrderPhaseAndStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Result:  "error",
			Message: "Invalid request",
			Content: err.Error(),
		})
		return
	}
	if err := h.service.WorkOrderPhaseIn(c.Request.Context(), req); err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			c.JSON(svcErr.StatusCode, models.ResponseMessage{
				Result:  "error",
				Message: svcErr.Message,
				Content: svcErr.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Result:  "error",
			Message: "Failed to set workorderphase and status",
			Content: err.Error(),
		})
		return
	}
	
	// Record metric
	observability.RecordWorkOrderPhaseIn(c.Request.Context(), req.WorkcenterID)
	
	c.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Workorderphase and status registered successfully",
	})
}

func(h *Handler)WorkOrderPhaseOut(c *gin.Context){
	var req models.WorkOrderPhaseOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Result:  "error",
			Message: "Invalid request",
			Content: err.Error(),
		})
		return
	}
	fmt.Println(req)
	if err := h.service.WorkOrderPhaseOut(c.Request.Context(), req); err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			c.JSON(svcErr.StatusCode, models.ResponseMessage{
				Result:  "error",
				Message: svcErr.Message,
				Content: svcErr.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Result:  "error",
			Message: "Failed to set workorderphase and status",
			Content: err.Error(),
		})
		return
	}
	
	// Record metric
	observability.RecordWorkOrderPhaseOut(c.Request.Context(), req.WorkcenterID)
	
	c.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Workorderphase and status registered successfully",
	})
}