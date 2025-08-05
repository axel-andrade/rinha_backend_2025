#!/bin/bash

# Script de benchmark para Rinha de Backend 2025
# Testa a performance das otimizações implementadas

echo "🚀 Iniciando benchmark da Rinha de Backend 2025"
echo "================================================"

# Configurações
BASE_URL="http://localhost:9999"
CONCURRENT_USERS=100
REQUESTS_PER_USER=100
TOTAL_REQUESTS=$((CONCURRENT_USERS * REQUESTS_PER_USER))

echo "📊 Configurações do teste:"
echo "   - URL Base: $BASE_URL"
echo "   - Usuários concorrentes: $CONCURRENT_USERS"
echo "   - Requests por usuário: $REQUESTS_PER_USER"
echo "   - Total de requests: $TOTAL_REQUESTS"
echo ""

# Função para gerar UUID
generate_uuid() {
    python3 -c "import uuid; print(str(uuid.uuid4()))"
}

# Função para gerar valor aleatório
generate_amount() {
    python3 -c "import random; print(round(random.uniform(1.0, 1000.0), 2))"
}

# Teste de health check
echo "🏥 Testando health check..."
health_response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/health")
if [ "$health_response" = "200" ]; then
    echo "   ✅ Health check OK"
else
    echo "   ❌ Health check falhou: $health_response"
    exit 1
fi

# Teste de performance - POST /payments
echo ""
echo "📈 Testando performance - POST /payments"
echo "   Iniciando $TOTAL_REQUESTS requests concorrentes..."

start_time=$(date +%s.%N)

# Criar arquivo temporário com requests
temp_file=$(mktemp)
for i in $(seq 1 $TOTAL_REQUESTS); do
    uuid=$(generate_uuid)
    amount=$(generate_amount)
    echo "curl -s -X POST $BASE_URL/payments -H 'Content-Type: application/json' -d '{\"correlationId\":\"$uuid\",\"amount\":$amount}'" >> "$temp_file"
done

# Executar requests em paralelo
parallel -j $CONCURRENT_USERS < "$temp_file" > /dev/null 2>&1

end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)

# Calcular métricas
requests_per_second=$(echo "scale=2; $TOTAL_REQUESTS / $duration" | bc)
avg_response_time=$(echo "scale=3; $duration * 1000 / $TOTAL_REQUESTS" | bc)

echo "   ✅ Teste concluído!"
echo "   📊 Métricas:"
echo "      - Duração: ${duration}s"
echo "      - Requests/segundo: ${requests_per_second}"
echo "      - Tempo médio de resposta: ${avg_response_time}ms"

# Teste de summary
echo ""
echo "📊 Testando GET /payments-summary..."
summary_response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/payments-summary")
if [ "$summary_response" = "200" ]; then
    echo "   ✅ Summary OK"
else
    echo "   ❌ Summary falhou: $summary_response"
fi

# Limpeza
rm -f "$temp_file"

echo ""
echo "🎯 Benchmark concluído!"
echo "================================================" 