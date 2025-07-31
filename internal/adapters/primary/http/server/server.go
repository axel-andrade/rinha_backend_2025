package http_server

import (
	"log"

	http_router "github.com/axel-andrade/rinha_backend_2025/internal/adapters/primary/http/router"
	"github.com/axel-andrade/rinha_backend_2025/internal/infra/bootstrap"
	"github.com/valyala/fasthttp"
)

type Server struct {
	handler fasthttp.RequestHandler
}

// NewServer já cria o handler com as rotas configuradas
func NewServer(d *bootstrap.Dependencies) *Server {
	return &Server{
		handler: http_router.ConfigureRoutes(d),
	}
}

// Run inicia o servidor na porta 9999
func (s *Server) Run() {
	log.Println("Server listening on :9999")
	if err := fasthttp.ListenAndServe(":9999", s.handler); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
