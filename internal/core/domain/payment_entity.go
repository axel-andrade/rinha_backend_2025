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
