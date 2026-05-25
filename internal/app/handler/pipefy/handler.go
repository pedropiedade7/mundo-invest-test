package pipefy

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pedropiedade7/mundo-invest-test/internal/app/handler/httphelper"
	"github.com/pedropiedade7/mundo-invest-test/internal/service"
)

type Handler struct {
	service *service.ProcessWebhookService
}

func New(service *service.ProcessWebhookService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ProcessCardUpdated(c *gin.Context) {
	var req cardUpdatedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.service.Execute(c.Request.Context(), service.ProcessWebhookInput{
		EventID:     req.EventID,
		CardID:      req.CardID,
		ClientEmail: req.ClientEmail,
		Timestamp:   req.Timestamp,
	})
	if err != nil {
		c.JSON(httphelper.StatusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

type cardUpdatedRequest struct {
	EventID     string    `json:"event_id" binding:"required"`
	CardID      string    `json:"card_id" binding:"required"`
	ClientEmail string    `json:"cliente_email" binding:"required,email"`
	Timestamp   time.Time `json:"timestamp" binding:"required"`
}
