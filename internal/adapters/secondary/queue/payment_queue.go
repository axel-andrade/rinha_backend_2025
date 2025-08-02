package queue

import (
	"context"
	"encoding/json"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentQueue struct {
	natsQueue  *NatsQueue
	topic      string
	queueGroup string
}

func NewPaymentQueue(natsQueue *NatsQueue) *PaymentQueue {
	return &PaymentQueue{
		natsQueue:  natsQueue,
		topic:      "payments",
		queueGroup: "payment-workers", // nome fixo para agrupar consumidores
	}
}

func (pq *PaymentQueue) PublishPayment(ctx context.Context, payment *domain.Payment) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return pq.natsQueue.Publish(pq.topic, payload)
}

func (pq *PaymentQueue) StartConsuming(ctx context.Context, handler func(context.Context, *domain.Payment) error) error {
	return pq.natsQueue.SubscribeQueue(pq.topic, pq.queueGroup, func(data []byte) {
		var payment domain.Payment
		if err := json.Unmarshal(data, &payment); err != nil {
			return
		}
		_ = handler(ctx, &payment)
	})
}

func (pq *PaymentQueue) StartConsumingWithWorkers(
	ctx context.Context,
	workerCount int,
	handler func(context.Context, *domain.Payment) error,
) error {
	return pq.natsQueue.SubscribeQueueWithWorkers(pq.topic, pq.queueGroup, func(data []byte) {
		var payment domain.Payment
		if err := json.Unmarshal(data, &payment); err != nil {
			return
		}

		// Chama handler com contexto
		if err := handler(ctx, &payment); err != nil {
		}
	}, workerCount)
}
