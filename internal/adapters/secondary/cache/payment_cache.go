package cache

import (
	"errors"
	"sync"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain/interfaces"
)

type InMemoryPaymentCache struct {
	healthCache sync.Map
	lockedKeys  map[string]struct{}
	mu          sync.Mutex
}

func NewInMemoryPaymentCache() *InMemoryPaymentCache {
	return &InMemoryPaymentCache{
		lockedKeys: make(map[string]struct{}),
	}
}

func (c *InMemoryPaymentCache) GetProcessorHealth(processor string) (interfaces.HealthStatus, error) {
	val, ok := c.healthCache.Load(processor)
	if !ok {
		return interfaces.HealthStatus{}, errors.New("not found")
	}
	return val.(interfaces.HealthStatus), nil
}

func (c *InMemoryPaymentCache) SetProcessorHealth(processor string, health interfaces.HealthStatus) error {
	c.healthCache.Store(processor, health)
	return nil
}

func (c *InMemoryPaymentCache) LockPayment(correlationId string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.lockedKeys[correlationId]; exists {
		return false, nil
	}

	c.lockedKeys[correlationId] = struct{}{}
	return true, nil
}

func (c *InMemoryPaymentCache) UnlockPayment(correlationId string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.lockedKeys, correlationId)
	return nil
}
