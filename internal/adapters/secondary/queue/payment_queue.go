package queue

import (
	"context"
	"encoding/json"

	"github.com/axel-andrade/rinha_backend_2025/internal/domain"
)

// PaymentQueue representa a fila de pagamentos que usa internamente NatsQueue.
type PaymentQueue struct {
	natsQueue *NatsQueue
	topic     string
}

// NewPaymentQueue cria uma instância de PaymentQueue usando uma NatsQueue.
func NewPaymentQueue(natsQueue *NatsQueue) *PaymentQueue {
	return &PaymentQueue{
		natsQueue: natsQueue,
		topic:     "payments",
	}
}

// PublishPayment publica um pagamento no tópico "payments".
func (pq *PaymentQueue) PublishPayment(ctx context.Context, payment *domain.Payment) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return pq.natsQueue.Publish(pq.topic, payload)
}
