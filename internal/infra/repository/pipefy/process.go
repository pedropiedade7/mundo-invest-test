package pipefy

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/pedropiedade7/mundo-invest-test/domain"
)

func (r *Repository) HasProcessedEvent(ctx context.Context, eventID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM eventos_processados WHERE event_id = $1)", eventID).Scan(&exists)
	return exists, err
}

func (r *Repository) MarkClientProcessed(ctx context.Context, eventID, cardID, email, priority string) (domain.Client, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Client{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "INSERT INTO eventos_processados (event_id, card_id) VALUES ($1, $2)", eventID, cardID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.Client{}, domain.ErrDuplicateEvent
		}
		return domain.Client{}, err
	}

	priorityID := domain.PriorityIDNormal
	if priority == domain.PriorityHigh {
		priorityID = domain.PriorityIDHigh
	}

	const query = `
		UPDATE clientes
		SET status_id = $1,
		    prioridade_id = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE email = $3
		RETURNING id, nome, email, tipo_solicitacao, valor_patrimonio, created_at, updated_at
	`

	var client domain.Client
	var createdAt, updatedAt time.Time

	err = tx.QueryRowContext(ctx, query, domain.StatusIDProcessed, priorityID, email).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.RequestType,
		&client.PortfolioCents,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return domain.Client{}, err
	}

	client.PortfolioValue = client.PortfolioCents / 100
	client.Status = domain.StatusProcessed
	client.Priority = &priority
	client.CreatedAt = createdAt.Format(time.RFC3339)
	client.UpdatedAt = updatedAt.Format(time.RFC3339)

	if err := tx.Commit(); err != nil {
		return domain.Client{}, err
	}
	return client, nil
}
