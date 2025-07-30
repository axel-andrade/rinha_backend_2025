package interfaces

import (
	"context"

	"github.com/axel-andrade/rinha_backend_2025/internal/core/domain"
)

type PaymentRepository interface {
	AddToStream(payment domain.Payment) error
	StorePayment(ctx context.Context, payment domain.Payment) error
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
}
