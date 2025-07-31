package bootstrap

import (
	http_adapter "github.com/axel-andrade/rinha_backend_2025/internal/adapters/primary/http"
	"github.com/axel-andrade/rinha_backend_2025/internal/adapters/secondary/database/postgres"
	"github.com/axel-andrade/rinha_backend_2025/internal/adapters/secondary/queue"
	"github.com/axel-andrade/rinha_backend_2025/internal/application"
	"github.com/axel-andrade/rinha_backend_2025/internal/domain/interfaces"
)

type Dependencies struct {
	HTTPHandler       http_adapter.Handler
	PaymentRepository interfaces.PaymentRepository
	// PaymentCache      application.PaymentCache // Add PaymentCache dependency
	PaymentService application.PaymentService
	SummaryService application.SummaryService
	PaymentQueue   interfaces.PaymentQueue
}

func LoadDependencies() *Dependencies {
	d := &Dependencies{}
	postgres.ConnectDB()
	db := postgres.GetDB()

	natsQueue := queue.NewNatsQueue()

	d.HTTPHandler = *http_adapter.NewHandler()
	d.PaymentRepository = postgres.NewPaymentRepository(db)
	// d.PaymentCache = application.NewPaymentCache() // Initialize PaymentCache
	d.PaymentQueue = queue.NewPaymentQueue(natsQueue)

	d.PaymentService = application.NewPaymentService(d.PaymentRepository, d.PaymentQueue)

	return d
}
