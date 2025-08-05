package payment_processor2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentProcessorClient struct {
	DefaultURL  string
	FallbackURL string
	HTTPClient  *http.Client
	// Connection pooling
	clientPool sync.Pool
}

func NewPaymentProcessorClient(defaultURL, fallbackURL string) *PaymentProcessorClient {
	// Configurações otimizadas para HTTP client
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true, // Desabilitar compressão para reduzir CPU
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   1 * time.Second, // Reduzir timeout para 1 segundo
	}

	return &PaymentProcessorClient{
		DefaultURL:  defaultURL,
		FallbackURL: fallbackURL,
		HTTPClient:  client,
		clientPool: sync.Pool{
			New: func() interface{} {
				return &http.Client{
					Transport: transport,
					Timeout:   1 * time.Second,
				}
			},
		},
	}
}

func (c *PaymentProcessorClient) ProcessPayment(ctx context.Context, payment domain.Payment) (domain.PaymentProcessorType, error) {
	if err := c.send(ctx, c.DefaultURL, payment); err == nil {
		return "default", nil
	}
	if err := c.send(ctx, c.FallbackURL, payment); err == nil {
		return "fallback", nil
	}
	return "", fmt.Errorf("failed to process payment in both processors")
}

func (c *PaymentProcessorClient) send(ctx context.Context, url string, payment domain.Payment) error {
	// Reutilizar client do pool
	client := c.clientPool.Get().(*http.Client)
	defer c.clientPool.Put(client)

	body := map[string]interface{}{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url+"/payments", bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("processor returned %d", resp.StatusCode)
}
