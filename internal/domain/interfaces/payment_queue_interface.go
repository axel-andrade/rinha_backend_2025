package interfaces

import (
	"context"

	"github.com/axel-andrade/rinha_backend_2025/internal/domain"
)

type PaymentQueue interface {
	PublishPayment(ctx context.Context, payment *domain.Payment) error
}
