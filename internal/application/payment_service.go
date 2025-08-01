package application

import (
	"context"
	"log"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain/interfaces"
)

type PaymentService struct {
	repository interfaces.PaymentRepository
	queue      interfaces.PaymentQueue
	processor  interfaces.PaymentProcessor
}

func NewPaymentService(
	repository interfaces.PaymentRepository,
	queue interfaces.PaymentQueue,
	processor interfaces.PaymentProcessor,
) *PaymentService {
	return &PaymentService{
		repository: repository,
		queue:      queue,
		processor:  processor,
	}
}

func (s *PaymentService) EnqueuePayment(ctx context.Context, p *domain.Payment) error {
	return s.queue.PublishPayment(ctx, p)
}

func (s *PaymentService) SavePayment(ctx context.Context, payment *domain.Payment) error {
	exists, err := s.repository.Exists(ctx, payment.CorrelationId.String())

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if err := s.repository.Save(ctx, *payment); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) GetPaymentSummary(ctx context.Context, from, to *time.Time) (domain.Summary, error) {
	return s.repository.GetSummary(ctx, from, to)
}

func (s *PaymentService) ProcessPayment(ctx context.Context, payment *domain.Payment) error {
	log.Printf("Processing payment: %s, Amount: %.2f", payment.CorrelationId, payment.Amount)

	processor, err := s.processor.ProcessPayment(ctx, *payment)
	if err != nil {
		log.Printf("Error processing payment %s: %v", payment.CorrelationId, err)
		// if enqueueErr := s.queue.PublishPayment(ctx, payment); enqueueErr != nil {
		// 	log.Printf("Failed to re-enqueue payment %s: %v", payment.CorrelationId, enqueueErr)
		// } else {
		// 	log.Printf("Payment %s re-enqueued for retry", payment.CorrelationId)
		// }
		return err
	}

	log.Printf("Payment %s processed successfully", processor)

	payment.Processor = processor

	if saveErr := s.SavePayment(ctx, payment); saveErr != nil {
		log.Printf("Failed to save payment %s: %v", payment.CorrelationId, saveErr)
	}

	return nil
}
