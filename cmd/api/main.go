package main

import (
	http_adapter "github.com/axel-andrade/rinha_backend_2025/internal/adapters/primary/http"
	"github.com/axel-andrade/rinha_backend_2025/internal/adapters/secondary/database/postgres"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	postgres.ConnectDB()
}

func main() {
	server := http_adapter.NewServer()
	server.Run()
}
