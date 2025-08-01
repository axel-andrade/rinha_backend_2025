package main

import (
	http_server "github.com/axel-andrade/go_rinha_backend_2025/internal/adapters/primary/http/server"
	"github.com/axel-andrade/go_rinha_backend_2025/internal/infra/bootstrap"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	d := bootstrap.LoadDependencies()
	server := http_server.NewServer(d)
	server.Run()
}
