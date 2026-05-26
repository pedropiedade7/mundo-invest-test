package tests

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/pedropiedade7/mundo-invest-test/internal/app/handler"
	clienthandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/client"
	pipefyhandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/pipefy"
	"github.com/pedropiedade7/mundo-invest-test/internal/domain"
	"github.com/pedropiedade7/mundo-invest-test/internal/infra/pipefy"
	"github.com/pedropiedade7/mundo-invest-test/internal/service"
)

func newTestRouter(repo *fakeRepo) http.Handler {
	pipefyClient := pipefy.NewClient("pipe_teste")
	createClient := service.NewCreateClientService(repo, pipefyClient)
	processWebhook := service.NewProcessWebhookService(repo, repo, pipefyClient)
	clientHandler := clienthandler.New(createClient)
	pipefyHandler := pipefyhandler.New(processWebhook)
	return handler.NewRouter(handler.NewHandler(clientHandler, pipefyHandler, nil))
}

type fakeRepo struct {
	mu                  sync.Mutex
	nextID              int64
	clients             map[string]domain.Client
	processedEvents     map[string]string
	webhookProcessCount int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		nextID:          1,
		clients:         make(map[string]domain.Client),
		processedEvents: make(map[string]string),
	}
}

func (r *fakeRepo) CreateClient(_ context.Context, client domain.Client) (domain.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.clients[client.Email]; ok {
		return domain.Client{}, domain.ErrClientAlreadyExists
	}

	client.ID = r.nextID
	r.nextID++
	if client.Status == "" {
		client.Status = domain.StatusPendingReview
	}
	if client.PortfolioCents == 0 {
		client.PortfolioCents = client.PortfolioValue * 100
	}
	if client.PortfolioValue == 0 {
		client.PortfolioValue = client.PortfolioCents / 100
	}
	r.clients[client.Email] = client
	return client, nil
}

func (r *fakeRepo) FindClientByEmail(_ context.Context, email string) (domain.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.clients[email]
	if !ok {
		return domain.Client{}, domain.ErrClientNotFound
	}
	return client, nil
}

func (r *fakeRepo) HasProcessedEvent(_ context.Context, eventID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.processedEvents[eventID]
	return ok, nil
}

func (r *fakeRepo) MarkClientProcessed(_ context.Context, eventID, cardID, email, priority string) (domain.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.processedEvents[eventID]; ok {
		return domain.Client{}, domain.ErrDuplicateEvent
	}

	client, ok := r.clients[email]
	if !ok {
		return domain.Client{}, domain.ErrClientNotFound
	}
	if priority == "" {
		return domain.Client{}, errors.New("prioridade vazia")
	}

	r.processedEvents[eventID] = cardID
	r.webhookProcessCount++
	client.Status = domain.StatusProcessed
	client.Priority = &priority
	r.clients[email] = client
	return client, nil
}
