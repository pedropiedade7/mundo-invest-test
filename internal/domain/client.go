package domain

import "errors"

const (
	StatusPendingReview = "Aguardando Análise"
	StatusProcessed     = "Processado"
	PriorityHigh        = "prioridade_alta"
	PriorityNormal      = "prioridade_normal"
)

const (
	StatusIDPendingReview int64 = 1
	StatusIDProcessed     int64 = 2

	PriorityIDNormal int64 = 1
	PriorityIDHigh   int64 = 2
)

var (
	ErrClientAlreadyExists = errors.New("cliente já cadastrado")
	ErrClientNotFound      = errors.New("cliente não encontrado")
	ErrDuplicateEvent      = errors.New("evento já processado")
)

type Client struct {
	ID             int64   `json:"id"`
	Name           string  `json:"cliente_nome"`
	Email          string  `json:"cliente_email"`
	RequestType    string  `json:"tipo_solicitacao"`
	PortfolioCents int64   `json:"-"`
	PortfolioValue int64   `json:"valor_patrimonio"`
	Status         string  `json:"status"`
	Priority       *string `json:"prioridade,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}
