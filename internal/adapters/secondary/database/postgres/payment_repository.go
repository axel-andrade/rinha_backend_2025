package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/axel-andrade/rinha_backend_2025/internal/domain"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) StorePayment(ctx context.Context, p domain.Payment) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO payments (id, amount, processor, requested_at)
		VALUES ($1, $2, $3, $4)
	`, p.CorrelationId, p.Amount, p.Processor, p.RequestedAt.UTC())

	if err != nil {
		return fmt.Errorf("failed to store payment: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetSummary(ctx context.Context, from, to *time.Time) (domain.Summary, error) {
	query := `
		SELECT processor, COUNT(*) AS total_requests, COALESCE(SUM(amount), 0) AS total_amount
		FROM payments
		WHERE ($1 IS NULL OR requested_at >= $1)
		  AND ($2 IS NULL OR requested_at <= $2)
		GROUP BY processor
	`

	rows, err := r.db.Query(ctx, query, from, to)
	if err != nil {
		return domain.Summary{}, fmt.Errorf("failed to query summary: %w", err)
	}
	defer rows.Close()

	// inicializa com zeros
	var summary domain.Summary

	for rows.Next() {
		var processor string
		var count int
		var amount float64

		if err := rows.Scan(&processor, &count, &amount); err != nil {
			return domain.Summary{}, fmt.Errorf("failed to scan summary row: %w", err)
		}

		switch processor {
		case string(domain.ProcessorDefault):
			summary.Default.TotalRequests = count
			summary.Default.TotalAmount = amount
		case string(domain.ProcessorFallback):
			summary.Fallback.TotalRequests = count
			summary.Fallback.TotalAmount = amount
		default:
			// opcional: logar processadores inesperados
		}
	}

	return summary, nil
}
