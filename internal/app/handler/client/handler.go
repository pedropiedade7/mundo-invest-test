package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedropiedade7/mundo-invest-test/internal/app/handler/httphelper"
	"github.com/pedropiedade7/mundo-invest-test/internal/service"
)

type Handler struct {
	service *service.CreateClientService
}

func New(service *service.CreateClientService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.service.Execute(c.Request.Context(), service.CreateClientInput{
		ClientName:     req.ClientName,
		ClientEmail:    req.ClientEmail,
		RequestType:    req.RequestType,
		PortfolioValue: req.PortfolioValue,
	})
	if err != nil {
		c.JSON(httphelper.StatusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

type createRequest struct {
	ClientName     string `json:"cliente_nome" binding:"required"`
	ClientEmail    string `json:"cliente_email" binding:"required,email"`
	RequestType    string `json:"tipo_solicitacao" binding:"required"`
	PortfolioValue int64  `json:"valor_patrimonio" binding:"required,gte=0"`
}
