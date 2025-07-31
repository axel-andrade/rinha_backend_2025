package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentProcessorType string

const (
	ProcessorDefault  PaymentProcessorType = "default"
	ProcessorFallback PaymentProcessorType = "fallback"
)

type Payment struct {
	CorrelationId uuid.UUID            `json:"correlationId"`
	Amount        float64              `json:"amount"`
	RequestedAt   time.Time            `json:"requestedAt"`
	Processor     PaymentProcessorType `json:"processor"`
}

func NewPayment(correlationId string, amount float64) *Payment {
	return &Payment{
		CorrelationId: uuid.MustParse(correlationId),
		Amount:        amount,
		RequestedAt:   time.Now().UTC(),
		Processor:     ProcessorDefault,
	}
}
