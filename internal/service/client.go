package service

import (
	"context"
	"strings"

	"github.com/pedropiedade7/mundo-invest-test/domain"
	"github.com/pedropiedade7/mundo-invest-test/internal/infra/pipefy"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, client domain.Client) (domain.Client, error)
}

type CreateClientInput struct {
	ClientName     string
	ClientEmail    string
	RequestType    string
	PortfolioValue int64
}

type CreateClientOutput struct {
	Client domain.Client `json:"cliente"`
}

type CreateClientService struct {
	repo   ClientRepository
	pipefy *pipefy.Client
}

func NewCreateClientService(repo ClientRepository, pipefyClient *pipefy.Client) *CreateClientService {
	return &CreateClientService{repo: repo, pipefy: pipefyClient}
}

func (s *CreateClientService) Execute(ctx context.Context, input CreateClientInput) (CreateClientOutput, error) {
	email := normalizeEmail(input.ClientEmail)
	if err := validateEmail(email); err != nil {
		return CreateClientOutput{}, err
	}

	client := domain.Client{
		Name:           strings.TrimSpace(input.ClientName),
		Email:          email,
		RequestType:    strings.TrimSpace(input.RequestType),
		PortfolioCents: input.PortfolioValue * 100,
		PortfolioValue: input.PortfolioValue,
		Status:         domain.StatusPendingReview,
	}

	saved, err := s.repo.CreateClient(ctx, client)
	if err != nil {
		return CreateClientOutput{}, err
	}

	// Simulates the Pipefy payload that would be sent to create the card.
	_ = s.pipefy.BuildCreateCardRequest(pipefy.CreateCardInput{
		ClientName:     saved.Name,
		ClientEmail:    saved.Email,
		RequestType:    saved.RequestType,
		PortfolioValue: saved.PortfolioValue,
	})

	return CreateClientOutput{Client: saved}, nil
}
