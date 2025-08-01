package interfaces

type PaymentCache interface {
	GetProcessorHealth(processor string) (HealthStatus, error)
	SetProcessorHealth(processor string, health HealthStatus) error
	LockPayment(correlationId string) (bool, error)
	UnlockPayment(correlationId string) error
}

type HealthStatus struct {
	Failing         bool
	MinResponseTime int
}
