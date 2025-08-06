# Rinha de Backend 2025 - Otimizações de Performance

## Otimizações Implementadas

### 1. **Banco de Dados (PostgreSQL)**
- ✅ Índices otimizados para consultas de summary
- ✅ Configurações de performance do PostgreSQL
- ✅ Query otimizada com índices compostos
- ✅ Configurações de autovacuum e WAL

### 2. **Servidor HTTP (FastHTTP)**
- ✅ Configurações de alta concorrência
- ✅ Object pooling para reduzir alocações
- ✅ Timeouts otimizados
- ✅ Buffer sizes otimizados
- ✅ Keep-alive habilitado

### 3. **Message Queue (NATS)**
- ✅ Buffer aumentado para 10k mensagens
- ✅ Workers dinâmicos baseados em CPU
- ✅ Configurações de reconexão otimizadas
- ✅ Acknowledgment para confiabilidade

### 4. **Payment Processor**
- ✅ Connection pooling
- ✅ Timeout reduzido para 1 segundo
- ✅ HTTP client otimizado
- ✅ Keep-alive habilitado

### 5. **Load Balancer (NGINX)**
- ✅ Worker connections aumentado para 8192
- ✅ Keep-alive otimizado
- ✅ Buffer configurations
- ✅ Gzip compression
- ✅ Proxy buffering

### 6. **Docker & Recursos**
- ✅ Recursos aumentados para todos os serviços
- ✅ Configurações de CPU e memória otimizadas
- ✅ PostgreSQL com configurações de performance

## Configurações de Performance

### Workers
- **NATS Workers**: 4x número de CPUs
- **FastHTTP Concurrency**: 1000x número de CPUs
- **NGINX Workers**: 8192 connections

### Recursos Docker
- **App Instances**: 1.0 CPU, 256MB RAM
- **PostgreSQL**: 0.5 CPU, 256MB RAM
- **NATS**: 0.25 CPU, 64MB RAM
- **NGINX**: 0.5 CPU, 64MB RAM

### Índices do Banco
- `idx_payments_processor_requested_at`
- `idx_payments_requested_at`
- `idx_payments_processor`

## Como Executar

```bash
docker-compose up --build
```

## Endpoints

- `POST /payments` - Enviar pagamento
- `GET /payments-summary` - Obter resumo
- `GET /health` - Health check

## Métricas Esperadas

Com essas otimizações, esperamos:
- ✅ Maior throughput de requests
- ✅ Menor latência
- ✅ Melhor utilização de recursos
- ✅ Maior taxa de sucesso nos processamentos