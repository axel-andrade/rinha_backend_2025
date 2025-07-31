package http_adapter

import (
	"log"

	"github.com/valyala/fasthttp"
)

type Server struct {
	handler fasthttp.RequestHandler
}

// NewServer já cria o handler com as rotas configuradas
func NewServer() *Server {
	h := NewHandler()
	return &Server{
		handler: ConfigureRoutes(h),
	}
}

// Run inicia o servidor na porta 9999
func (s *Server) Run() {
	log.Println("Server listening on :9999")
	if err := fasthttp.ListenAndServe(":9999", s.handler); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
