package service

import (
	"context"
	"time"

	"github.com/pedropiedade7/mundo-invest-test/internal/domain"
	"github.com/pedropiedade7/mundo-invest-test/internal/infra/pipefy"
)

type ClientFinder interface {
	FindClientByEmail(ctx context.Context, email string) (domain.Client, error)
}

type PipefyEventRepository interface {
	HasProcessedEvent(ctx context.Context, eventID string) (bool, error)
	MarkClientProcessed(ctx context.Context, eventID, cardID, email, priority string) (domain.Client, error)
}

type ProcessWebhookInput struct {
	EventID     string
	CardID      string
	ClientEmail string
	Timestamp   time.Time
}

type ProcessWebhookOutput struct {
	Client domain.Client `json:"cliente"`
}

type ProcessWebhookService struct {
	clientRepo ClientFinder
	pipefyRepo PipefyEventRepository
	pipefy     *pipefy.Client
}

func NewProcessWebhookService(clientRepo ClientFinder, pipefyRepo PipefyEventRepository, pipefyClient *pipefy.Client) *ProcessWebhookService {
	return &ProcessWebhookService{
		clientRepo: clientRepo,
		pipefyRepo: pipefyRepo,
		pipefy:     pipefyClient,
	}
}

func (s *ProcessWebhookService) Execute(ctx context.Context, input ProcessWebhookInput) (ProcessWebhookOutput, error) {
	email := normalizeEmail(input.ClientEmail)
	if err := validateEmail(email); err != nil {
		return ProcessWebhookOutput{}, err
	}

	processed, err := s.pipefyRepo.HasProcessedEvent(ctx, input.EventID)
	if err != nil {
		return ProcessWebhookOutput{}, err
	}
	if processed {
		return ProcessWebhookOutput{}, domain.ErrDuplicateEvent
	}

	client, err := s.clientRepo.FindClientByEmail(ctx, email)
	if err != nil {
		return ProcessWebhookOutput{}, err
	}

	priority := domain.PriorityNormal

	if client.PortfolioCents >= 20000000 {
		priority = domain.PriorityHigh
	}

	updated, err := s.pipefyRepo.MarkClientProcessed(ctx, input.EventID, input.CardID, email, priority)
	if err != nil {
		return ProcessWebhookOutput{}, err
	}

	// Simulates the Pipefy payload that would update status and priority fields.
	_ = s.pipefy.BuildUpdateCardFieldsRequest(input.CardID, domain.StatusProcessed, priority)

	return ProcessWebhookOutput{Client: updated}, nil
}
