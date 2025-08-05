package bootstrap

import (
	"context"
	"runtime"

	http_handler "github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/primary/http/handler"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/secondary/database/postgres"
	payment_processor2 "github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/secondary/payment_processor"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/secondary/queue"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/application"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/domain/interfaces"
)

type Dependencies struct {
	PaymentRepository interfaces.PaymentRepository
	PaymentService    *application.PaymentService
	SummaryService    application.SummaryService
	PaymentQueue      interfaces.PaymentQueue
	PaymentProcessor  interfaces.PaymentProcessor
	HTTPHandler       *http_handler.Handler
}

func LoadDependencies() *Dependencies {
	d := &Dependencies{}
	postgres.ConnectDB()
	db := postgres.GetDB()
	natsQueue := queue.NewNatsQueue()
	paymentQueue := queue.NewPaymentQueue(natsQueue)
	paymentProcessor := payment_processor2.NewPaymentProcessorClient("http://payment-processor-default:8080", "http://payment-processor-fallback:8080")

	d.PaymentRepository = postgres.NewPaymentRepository(db)
	d.PaymentQueue = paymentQueue
	d.PaymentProcessor = paymentProcessor

	pService := application.NewPaymentService(d.PaymentRepository, d.PaymentQueue, d.PaymentProcessor)

	d.PaymentService = pService

	// Aumentar workers para melhor performance
	workerCount := runtime.NumCPU() * 4 // 4x o número de CPUs
	paymentQueue.StartConsumingWithWorkers(context.Background(), workerCount, pService.ProcessPayment)

	d.HTTPHandler = http_handler.NewHandler(pService)

	return d
}
