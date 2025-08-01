package interfaces

import (
	"context"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, p domain.Payment) (domain.PaymentProcessorType, error)
}
