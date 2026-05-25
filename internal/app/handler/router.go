package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clienthandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/client"
	pipefyhandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/pipefy"
)

type HealthChecker interface {
	Ping() error
}

type Handler struct {
	client *clienthandler.Handler
	pipefy *pipefyhandler.Handler
	db     HealthChecker
}

func NewHandler(client *clienthandler.Handler, pipefy *pipefyhandler.Handler, db HealthChecker) *Handler {
	return &Handler{
		client: client,
		pipefy: pipefy,
		db:     db,
	}
}

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health)
	r.POST("/clientes", handler.client.Create)
	r.POST("/webhooks/pipefy/card-updated", handler.pipefy.ProcessCardUpdated)

	return r
}

func (h *Handler) Health(c *gin.Context) {
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "UP", "database": "DOWN: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "UP", "database": "CONNECTED"})
}
