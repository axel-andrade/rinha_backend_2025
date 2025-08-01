package payment_processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentProcessorClient struct {
	DefaultURL  string
	FallbackURL string
	HTTPClient  *http.Client
	cacheMutex  sync.Mutex

	defaultFailUntil  time.Time
	defaultRetries    int
	maxDefaultRetries int
}

func NewPaymentProcessorClient(defaultURL, fallbackURL string) *PaymentProcessorClient {
	return &PaymentProcessorClient{
		DefaultURL:        defaultURL,
		FallbackURL:       fallbackURL,
		HTTPClient:        &http.Client{Timeout: 2 * time.Second},
		maxDefaultRetries: 3,
	}
}

func (c *PaymentProcessorClient) ProcessPayment(ctx context.Context, payment domain.Payment) (domain.PaymentProcessorType, error) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	now := time.Now()
	skipDefault := now.Before(c.defaultFailUntil) || c.defaultRetries >= c.maxDefaultRetries

	if !skipDefault {
		if err := c.send(ctx, c.DefaultURL, payment); err == nil {
			c.defaultRetries = 0
			c.defaultFailUntil = time.Time{}
			return domain.ProcessorDefault, nil
		} else {
			log.Printf("[Default Processor] attempt %d failed: %v", c.defaultRetries+1, err)
			c.defaultRetries++
			if c.defaultRetries >= c.maxDefaultRetries {
				c.defaultFailUntil = now.Add(5 * time.Second)
				log.Printf("[Default Processor] entering cooldown until %s", c.defaultFailUntil.Format(time.RFC3339))
			}
		}
	}

	if err := c.send(ctx, c.FallbackURL, payment); err == nil {
		c.defaultRetries = 0
		c.defaultFailUntil = time.Time{}
		return domain.ProcessorFallback, nil
	} else {
		log.Printf("[Fallback Processor] attempt failed: %v", err)
	}

	c.defaultRetries = 0
	c.defaultFailUntil = time.Time{}

	return "", fmt.Errorf("failed to process payment in both processors")
}

func (c *PaymentProcessorClient) send(ctx context.Context, url string, payment domain.Payment) error {
	body := map[string]interface{}{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/payments", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	bodyStr := buf.String()
	if len(bodyStr) > 500 {
		bodyStr = bodyStr[:500] + "...(truncated)"
	}

	return fmt.Errorf("processor returned status %d: %s", resp.StatusCode, bodyStr)
}
