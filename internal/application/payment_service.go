package application

import (
	"context"
	"time"

	"github.com/axel-andrade/rinha_backend_2025/internal/domain"
	"github.com/axel-andrade/rinha_backend_2025/internal/domain/interfaces"
)

type PaymentService struct {
	repository interfaces.PaymentRepository
	// cache      interfaces.PaymentCache
	queue interfaces.PaymentQueue
}

func NewPaymentService(
	repository interfaces.PaymentRepository,
	// cache interfaces.PaymentCache,
	queue interfaces.PaymentQueue,
) *PaymentService {
	return &PaymentService{
		repository: repository,
		// cache:      cache,
		queue: queue,
	}
}

// Enfileira o pagamento com correlationId e valor
func (s *PaymentService) EnqueuePayment(ctx context.Context, correlationId string, amount float64) error {
	payment := domain.NewPayment(correlationId, amount)
	return s.queue.PublishPayment(ctx, payment)
}

// Salva o pagamento se ainda não existe (idempotente) e opcionalmente armazena no cache
func (s *PaymentService) SavePayment(ctx context.Context, payment *domain.Payment) error {
	// exists, err := s.repository.Exists(ctx, payment.CorrelationId.String())
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return nil // já processado
	// }

	if err := s.repository.Save(ctx, *payment); err != nil {
		return err
	}

	// if err := s.cache.SetPayment(ctx, payment); err != nil {
	// 	return err
	// }

	return nil
}

// Recupera um resumo de pagamentos entre datas
func (s *PaymentService) GetPaymentSummary(ctx context.Context, from, to *time.Time) (domain.Summary, error) {
	return s.repository.GetSummary(ctx, from, to)
}

// Processa o pagamento vindo da fila (pode incluir persistência ou lógica adicional no futuro)
func (s *PaymentService) ProcessPayment(ctx context.Context, payment *domain.Payment) error {
	// Lógica futura de validação, persistência ou chamada externa
	// No momento, você pode apenas fazer um log ou contagem
	return nil
}
