package interfaces

import (
	"context"
	"time"

	"github.com/axel-andrade/rinha_backend_2025/internal/domain"
)

// PaymentRepository define as operações para persistência e consulta de pagamentos
type PaymentRepository interface {
	// Save registra um pagamento, retorna erro se falhar
	Save(ctx context.Context, payment domain.Payment) error

	// Exists verifica se já existe um pagamento com o correlationId (idempotência)
	Exists(ctx context.Context, correlationId string) (bool, error)

	// GetSummary retorna o resumo de pagamentos entre o intervalo opcional de datas
	GetSummary(ctx context.Context, from, to *time.Time) (domain.Summary, error)

	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
}
