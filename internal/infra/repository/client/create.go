package client

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/pedropiedade7/mundo-invest-test/domain"
)

func (r *Repository) CreateClient(ctx context.Context, client domain.Client) (domain.Client, error) {
	const query = `
		INSERT INTO clientes (nome, email, tipo_solicitacao, valor_patrimonio, status_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, nome, email, tipo_solicitacao, valor_patrimonio, created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		client.Name,
		client.Email,
		client.RequestType,
		client.PortfolioCents,
		domain.StatusIDPendingReview,
	).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.RequestType,
		&client.PortfolioCents,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.Client{}, domain.ErrClientAlreadyExists
		}
		return domain.Client{}, err
	}

	client.Status = domain.StatusPendingReview
	client.PortfolioValue = client.PortfolioCents / 100
	client.CreatedAt = createdAt.Format(time.RFC3339)
	client.UpdatedAt = updatedAt.Format(time.RFC3339)
	return client, nil
}
