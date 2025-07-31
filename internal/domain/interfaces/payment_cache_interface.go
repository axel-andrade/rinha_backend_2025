package interfaces

type PaymentCache interface {
	// GetProcessorHealth retorna o status cached do processor ("default" ou "fallback")
	GetProcessorHealth(processor string) (HealthStatus, error)

	// SetProcessorHealth atualiza o cache do health
	SetProcessorHealth(processor string, health HealthStatus) error

	// LockPayment tenta adquirir um lock para processar um correlationId específico
	LockPayment(correlationId string) (bool, error)

	// UnlockPayment libera o lock
	UnlockPayment(correlationId string) error
}

type HealthStatus struct {
	Failing         bool
	MinResponseTime int
}
