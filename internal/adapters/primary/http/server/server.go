package http_server

import (
	"log"
	"runtime"

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
	// Configurações otimizadas para performance
	server := &fasthttp.Server{
		Handler:                      s.handler,
		ReadTimeout:                  5 * 1000,    // 5 segundos
		WriteTimeout:                 5 * 1000,    // 5 segundos
		IdleTimeout:                  10 * 1000,   // 10 segundos
		MaxRequestBodySize:           1024 * 1024, // 1MB
		DisablePreParseMultipartForm: true,
		NoDefaultContentType:         true,
		NoDefaultServerHeader:        true,
		// Configurações de concorrência
		Concurrency: runtime.NumCPU() * 1000, // Muito alto para alta concorrência
		// Configurações de buffer
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		// Configurações de keep-alive
		GetOnly:           false,
		DisableKeepalive:  false,
		KeepHijackedConns: true,
		// Configurações adicionais para resolver problema de conexão
		CloseOnShutdown: true,
		Logger:          nil, // Desabilitar logging interno
	}

	log.Println("Server listening on :9999")
	if err := server.ListenAndServe(":9999"); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
