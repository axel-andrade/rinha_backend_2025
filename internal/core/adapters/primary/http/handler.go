package http

import (
	"encoding/json"
	"log"
	"time"

	"github.com/axel-andrade/rinha_backend_2025/internal/core/domain"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	// PaymentSvc *application.PaymentService
	// SummarySvc *application.SummaryService
}

type paymentRequest struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

// POST /payments
func (h *Handler) HandlePayments(ctx *fasthttp.RequestCtx) {
	var req paymentRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		log.Printf("[Handler] Invalid body: %v", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid body")
		return
	}

	id, err := uuid.Parse(req.CorrelationId)
	if err != nil {
		log.Printf("[Handler] Invalid correlationId: %v", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid correlationId")
		return
	}

	if req.Amount <= 0 {
		log.Printf("[Handler] Invalid amount: %v", req.Amount)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount must be > 0")
		return
	}

	payment := domain.Payment{
		CorrelationId: id,
		Amount:        req.Amount,
		RequestedAt:   time.Now().UTC(),
	}

	data, _ := json.Marshal(payment)
	err = h.PaymentSvc.RepoImpl.Client.LPush(ctx, "payments:queue", data).Err()
	if err != nil {
		log.Printf("[Handler] Failed to enqueue payment: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to enqueue payment")
		return
	}

	log.Printf("[Handler] Payment enqueued: %s", req.CorrelationId)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	ctx.SetContentType("application/json")
	resp := map[string]string{
		"status":        "success",
		"message":       "Payment request accepted",
		"correlationId": req.CorrelationId,
	}
	respBytes, _ := json.Marshal(resp)
	ctx.SetBody(respBytes)
}

// GET /payments-summary?from=...&to=...
func (h *Handler) HandleSummary(ctx *fasthttp.RequestCtx) {
	fromStr := string(ctx.QueryArgs().Peek("from"))
	toStr := string(ctx.QueryArgs().Peek("to"))

	var from, to *time.Time
	if fromStr != "" {
		if f, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &f
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	summary, err := h.SummarySvc.GetSummary(from, to)
	if err != nil {
		log.Printf("[Handler] Failed to get summary: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to get summary")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	respBytes, err := json.Marshal(summary)
	if err != nil {
		log.Printf("[Handler] Failed to marshal summary: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to marshal summary")
		return
	}
	ctx.SetBody(respBytes)
}

// GET /health
func (h *Handler) HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("ok")
}
