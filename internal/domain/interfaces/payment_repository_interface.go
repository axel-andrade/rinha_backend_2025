package interfaces

import (
	"context"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, p domain.Payment) error
	Exists(ctx context.Context, correlationId string) (bool, error)
	GetSummary(ctx context.Context, from, to *time.Time) (domain.Summary, error)
}
