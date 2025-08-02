package payment_processor2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
)

type PaymentProcessorClient struct {
	DefaultURL  string
	FallbackURL string
	HTTPClient  *http.Client
}

func NewPaymentProcessorClient(defaultURL, fallbackURL string) *PaymentProcessorClient {
	return &PaymentProcessorClient{
		DefaultURL:  defaultURL,
		FallbackURL: fallbackURL,
		HTTPClient:  &http.Client{Timeout: 2 * time.Second},
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
	body := map[string]interface{}{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", url+"/payments", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("processor returned %d", resp.StatusCode)
}
