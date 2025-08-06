package http_handler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/axel-andrade/go_rinha_backend_2025/internal/application"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	PaymentService *application.PaymentService
	paymentPool    sync.Pool
}

func NewHandler(p *application.PaymentService) *Handler {
	return &Handler{
		PaymentService: p,
		paymentPool: sync.Pool{
			New: func() interface{} {
				return &paymentRequest{}
			},
		},
	}
}

type paymentRequest struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func (h *Handler) HandlePayments(ctx *fasthttp.RequestCtx) {
	req := h.paymentPool.Get().(*paymentRequest)
	defer h.paymentPool.Put(req)

	req.CorrelationId = ""
	req.Amount = 0

	if err := json.Unmarshal(ctx.PostBody(), req); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid body")
		return
	}

	id, err := uuid.Parse(req.CorrelationId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid correlationId")
		return
	}

	if req.Amount <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount must be > 0")
		return
	}

	payment := domain.Payment{
		CorrelationId: id,
		Amount:        req.Amount,
		RequestedAt:   time.Now().UTC(),
	}

	if err := h.PaymentService.EnqueuePayment(context.Background(), &payment); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to enqueue payment")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	ctx.SetBodyString("")
}

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

	summary, err := h.PaymentService.GetPaymentSummary(context.Background(), from, to)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to get payment summary")
		return
	}

	respBytes, err := json.Marshal(summary)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("failed to marshal summary")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(respBytes)
}

func (h *Handler) HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("ok")
}
