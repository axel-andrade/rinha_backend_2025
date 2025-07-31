package bootstrap

import (
	"context"

	http_handler "github.com/axel-andrade/rinha_backend_2025/internal/adapters/primary/http/handler"
	"github.com/axel-andrade/rinha_backend_2025/internal/adapters/secondary/database/postgres"
	"github.com/axel-andrade/rinha_backend_2025/internal/adapters/secondary/queue"
	"github.com/axel-andrade/rinha_backend_2025/internal/application"
	"github.com/axel-andrade/rinha_backend_2025/internal/domain/interfaces"
)

type Dependencies struct {
	PaymentRepository interfaces.PaymentRepository
	// PaymentCache      application.PaymentCache // Add PaymentCache dependency
	PaymentService application.PaymentService
	SummaryService application.SummaryService
	PaymentQueue   interfaces.PaymentQueue
	HTTPHandler    http_handler.Handler
}

func LoadDependencies() *Dependencies {
	d := &Dependencies{}
	postgres.ConnectDB()
	db := postgres.GetDB()
	natsQueue := queue.NewNatsQueue()
	paymentQueue := queue.NewPaymentQueue(natsQueue)

	d.PaymentRepository = postgres.NewPaymentRepository(db)
	d.PaymentQueue = paymentQueue
	pService := *application.NewPaymentService(d.PaymentRepository, d.PaymentQueue)

	d.PaymentService = pService
	paymentQueue.StartConsuming(context.Background(), pService.ProcessPayment)

	d.HTTPHandler = *http_handler.NewHandler(&d.PaymentService)

	return d
}
