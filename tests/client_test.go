package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pedropiedade7/mundo-invest-test/internal/domain"
)

func TestCreateClientWithValidPayloadPersistsClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	router := newTestRouter(repo)

	body := `{
		"cliente_nome": "João Silva",
		"cliente_email": "joao.silva@example.com",
		"tipo_solicitacao": "Atualização cadastral",
		"valor_patrimonio": 250000
	}`

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	client, err := repo.FindClientByEmail(context.Background(), "joao.silva@example.com")
	if err != nil {
		t.Fatalf("cliente não foi salvo: %v", err)
	}
	if client.Status != domain.StatusPendingReview {
		t.Fatalf("status = %q, want %q", client.Status, domain.StatusPendingReview)
	}

	if strings.Contains(rec.Body.String(), "pipefy_payload") {
		t.Fatalf("pipefy_payload não deve ser exposto na API: %s", rec.Body.String())
	}
}

func TestCreateClientWithDuplicatedEmailReturnsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newFakeRepo()
	router := newTestRouter(repo)

	body := `{
		"cliente_nome": "João Silva",
		"cliente_email": "joao.silva@example.com",
		"tipo_solicitacao": "Atualização cadastral",
		"valor_patrimonio": 250000
	}`

	first := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(first, req)
	if first.Code != http.StatusCreated {
		t.Fatalf("primeiro status = %d, body = %s", first.Code, first.Body.String())
	}

	second := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(second, req)
	if second.Code != http.StatusConflict {
		t.Fatalf("segundo status = %d, want %d, body = %s", second.Code, http.StatusConflict, second.Body.String())
	}
}
