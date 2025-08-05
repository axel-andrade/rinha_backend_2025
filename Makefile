.PHONY: build up down logs clean benchmark test health

# Configurações
DOCKER_COMPOSE = docker-compose
SERVICE_NAME = app1

# Comandos principais
build:
	@echo "🔨 Construindo containers..."
	$(DOCKER_COMPOSE) build

up:
	@echo "🚀 Iniciando serviços..."
	$(DOCKER_COMPOSE) up -d

down:
	@echo "🛑 Parando serviços..."
	$(DOCKER_COMPOSE) down

restart:
	@echo "🔄 Reiniciando serviços..."
	$(DOCKER_COMPOSE) restart

logs:
	@echo "📋 Exibindo logs..."
	$(DOCKER_COMPOSE) logs -f

# Comandos de desenvolvimento
dev: build up
	@echo "🎯 Ambiente de desenvolvimento iniciado!"
	@echo "📊 Acesse: http://localhost:9999"
	@echo "📋 Logs: make logs"

clean:
	@echo "🧹 Limpando containers e volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f

# Comandos de teste e benchmark
test:
	@echo "🧪 Executando testes..."
	curl -f http://localhost:9999/health || (echo "❌ Serviço não está respondendo" && exit 1)
	@echo "✅ Health check passou!"

benchmark:
	@echo "📈 Executando benchmark..."
	@if [ ! -f benchmark.sh ]; then echo "❌ Arquivo benchmark.sh não encontrado"; exit 1; fi
	@./benchmark.sh

health:
	@echo "🏥 Verificando saúde dos serviços..."
	@echo "NGINX:"
	@curl -s -w "Status: %{http_code}, Tempo: %{time_total}s\n" -o /dev/null http://localhost:9999/health || echo "❌ NGINX não responde"
	@echo "App1:"
	@curl -s -w "Status: %{http_code}, Tempo: %{time_total}s\n" -o /dev/null http://localhost:8080/health || echo "❌ App1 não responde"
	@echo "PostgreSQL:"
	@docker exec postgres-db pg_isready -U postgres || echo "❌ PostgreSQL não responde"
	@echo "NATS:"
	@docker exec nats nats-server --version || echo "❌ NATS não responde"

# Comandos de monitoramento
monitor:
	@echo "📊 Monitorando recursos..."
	@echo "CPU e Memória:"
	@docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
	@echo ""
	@echo "Logs de erro:"
	@$(DOCKER_COMPOSE) logs --tail=10 | grep -i error || echo "✅ Nenhum erro encontrado"

# Comandos de performance
perf-test:
	@echo "⚡ Teste de performance rápido..."
	@for i in {1..10}; do \
		curl -s -X POST http://localhost:9999/payments \
			-H "Content-Type: application/json" \
			-d '{"correlationId":"123e4567-e89b-12d3-a456-426614174000","amount":100.50}' \
			-w "Request $$i: %{http_code} - %{time_total}s\n" -o /dev/null; \
	done

# Help
help:
	@echo "🎯 Rinha de Backend 2025 - Comandos disponíveis:"
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  make build     - Construir containers"
	@echo "  make up        - Iniciar serviços"
	@echo "  make down      - Parar serviços"
	@echo "  make restart   - Reiniciar serviços"
	@echo "  make dev       - Iniciar ambiente de desenvolvimento"
	@echo "  make logs      - Exibir logs"
	@echo "  make clean     - Limpar tudo"
	@echo ""
	@echo "🧪 Testes:"
	@echo "  make test      - Teste básico"
	@echo "  make health    - Verificar saúde dos serviços"
	@echo "  make benchmark - Executar benchmark completo"
	@echo "  make perf-test - Teste de performance rápido"
	@echo "  make monitor   - Monitorar recursos"
	@echo ""
	@echo "📊 Endpoints:"
	@echo "  POST /payments         - Enviar pagamento"
	@echo "  GET  /payments-summary - Obter resumo"
	@echo "  GET  /health           - Health check"