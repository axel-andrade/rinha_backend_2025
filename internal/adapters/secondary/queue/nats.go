package queue

import (
	"log"
	"os"
	"sync"

	"github.com/nats-io/nats.go"
)

// NatsQueue representa uma conexão genérica ao NATS.
type NatsQueue struct {
	conn *nats.Conn
}

var (
	natsOnce sync.Once
	natsInst *NatsQueue
)

// NewNatsQueue cria ou retorna uma instância singleton da conexão NATS.
func NewNatsQueue() *NatsQueue {
	natsOnce.Do(func() {
		natsURL := os.Getenv("NATS_URL")
		if natsURL == "" {
			natsURL = "nats://localhost:4222"
		}

		conn, err := nats.Connect(natsURL)
		if err != nil {
			log.Fatalf("Erro ao conectar ao NATS: %v", err)
		}

		natsInst = &NatsQueue{
			conn: conn,
		}
	})
	return natsInst
}

// Publish publica uma mensagem no tópico informado.
func (q *NatsQueue) Publish(topic string, data []byte) error {
	return q.conn.Publish(topic, data)
}
