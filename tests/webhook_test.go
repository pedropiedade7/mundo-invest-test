package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pedropiedade7/mundo-invest-test/domain"
)

func TestWebhookAppliesHighPriorityByPortfolioValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	_, err := repo.CreateClient(context.Background(), domain.Client{
		Name:           "João Silva",
		Email:          "joao.silva@example.com",
		RequestType:    "Atualização cadastral",
		PortfolioValue: 250000,
		Status:         domain.StatusPendingReview,
	})
	if err != nil {
		t.Fatal(err)
	}

	router := newTestRouter(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(`{
		"event_id": "evt_123",
		"card_id": "card_456",
		"cliente_email": "joao.silva@example.com",
		"timestamp": "2026-05-18T12:00:00Z"
	}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	client, err := repo.FindClientByEmail(context.Background(), "joao.silva@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if client.Status != domain.StatusProcessed {
		t.Fatalf("status = %q, want %q", client.Status, domain.StatusProcessed)
	}
	if client.Priority == nil || *client.Priority != domain.PriorityHigh {
		t.Fatalf("prioridade = %v, want %q", client.Priority, domain.PriorityHigh)
	}

	if strings.Contains(rec.Body.String(), "pipefy_payload") {
		t.Fatalf("pipefy_payload não deve ser exposto na API: %s", rec.Body.String())
	}
}

func TestDuplicateWebhookIsNotProcessedAgain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	_, err := repo.CreateClient(context.Background(), domain.Client{
		Name:           "Maria Souza",
		Email:          "maria.souza@example.com",
		RequestType:    "Abertura de conta",
		PortfolioValue: 150000,
		Status:         domain.StatusPendingReview,
	})
	if err != nil {
		t.Fatal(err)
	}

	router := newTestRouter(repo)
	payload := `{
		"event_id": "evt_dup",
		"card_id": "card_999",
		"cliente_email": "maria.souza@example.com",
		"timestamp": "2026-05-18T12:00:00Z"
	}`

	first := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(first, req)
	if first.Code != http.StatusOK {
		t.Fatalf("primeiro status = %d, body = %s", first.Code, first.Body.String())
	}

	second := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(second, req)
	if second.Code != http.StatusConflict {
		t.Fatalf("segundo status = %d, want %d, body = %s", second.Code, http.StatusConflict, second.Body.String())
	}
	if repo.webhookProcessCount != 1 {
		t.Fatalf("processamentos = %d, want 1", repo.webhookProcessCount)
	}
}
