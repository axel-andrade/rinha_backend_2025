package http_server

import (
	"log"
	"runtime"
	"time"

	http_router "github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/primary/http/router"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/infra/bootstrap"
	"github.com/valyala/fasthttp"
)

type Server struct {
	handler fasthttp.RequestHandler
}

func NewServer(d *bootstrap.Dependencies) *Server {
	return &Server{
		handler: http_router.ConfigureRoutes(d),
	}
}

func (s *Server) Run() {
	server := &fasthttp.Server{
		Handler:                      s.handler,
		ReadTimeout:                  5 * time.Second,
		WriteTimeout:                 5 * time.Second,
		IdleTimeout:                  10 * time.Second,
		MaxRequestBodySize:           1024 * 1024, // 1MB
		DisablePreParseMultipartForm: true,
		NoDefaultContentType:         true,
		NoDefaultServerHeader:        true,
		Concurrency:                  runtime.NumCPU() * 100,
		ReadBufferSize:               4096,
		WriteBufferSize:              4096,
		DisableKeepalive:             false,
		CloseOnShutdown:              true,
		Logger:                       nil,
	}

	log.Println("Server listening on :9999")
	if err := server.ListenAndServe(":9999"); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
