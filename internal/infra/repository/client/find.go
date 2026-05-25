package client

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pedropiedade7/mundo-invest-test/domain"
)

func (r *Repository) FindClientByEmail(ctx context.Context, email string) (domain.Client, error) {
	const query = `
		SELECT c.id, c.nome, c.email, c.tipo_solicitacao, c.valor_patrimonio,
		       s.descricao, p.descricao, c.created_at, c.updated_at
		FROM clientes c
		JOIN cliente_status s ON s.id = c.status_id
		LEFT JOIN cliente_prioridades p ON p.id = c.prioridade_id
		WHERE c.email = $1
	`

	var client domain.Client
	var priority sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.RequestType,
		&client.PortfolioCents,
		&client.Status,
		&priority,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Client{}, domain.ErrClientNotFound
	}
	if err != nil {
		return domain.Client{}, err
	}

	client.PortfolioValue = client.PortfolioCents / 100
	if priority.Valid {
		client.Priority = &priority.String
	}
	client.CreatedAt = createdAt.Format(time.RFC3339)
	client.UpdatedAt = updatedAt.Format(time.RFC3339)
	return client, nil
}
