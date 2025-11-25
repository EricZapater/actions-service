package status

import (
	"actions-service/internal/models"
	"errors"
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

func (h *Handler) StatusIn(c *gin.Context) {
	var req models.StatusDTORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Result:  "error",
			Message: "Invalid request",
			Content: err.Error(),
		})
		return
	}	
	var statusReasonId *string
	if req.StatusReasonId != nil {
		reasonStr := req.StatusReasonId.String()
		statusReasonId = &reasonStr
	}	
	if err := h.service.StatusIn(c.Request.Context(), req.WorkcenterID.String(), req.StatusID.String(), statusReasonId); err != nil {
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
			Message: "Failed to set status in",
			Content: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.ResponseMessage{
		Result:  "success",
		Message: "Status in registered successfully",
	})
}