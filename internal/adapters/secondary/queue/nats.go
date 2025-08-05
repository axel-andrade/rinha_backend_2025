package queue

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsQueue struct {
	conn *nats.Conn
}

var (
	natsOnce sync.Once
	natsInst *NatsQueue
)

func NewNatsQueue() *NatsQueue {
	natsOnce.Do(func() {
		natsURL := os.Getenv("NATS_URL")
		if natsURL == "" {
			natsURL = "nats://localhost:4222"
		}

		// Configurações otimizadas para NATS
		opts := []nats.Option{
			nats.Name("rinha-backend"),
			nats.ReconnectWait(1 * time.Second),
			nats.MaxReconnects(10),
			nats.ReconnectJitter(100*time.Millisecond, 1*time.Second),
			nats.Timeout(5 * time.Second),
			nats.FlusherTimeout(5 * time.Second),
		}

		conn, err := nats.Connect(natsURL, opts...)
		if err != nil {
			log.Fatalf("Erro ao conectar ao NATS: %v", err)
		}

		natsInst = &NatsQueue{
			conn: conn,
		}
	})
	return natsInst
}

func (q *NatsQueue) Publish(topic string, data []byte) error {
	return q.conn.Publish(topic, data)
}

func (q *NatsQueue) SubscribeQueue(topic, queueGroup string, handler func(data []byte)) error {
	_, err := q.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {
		handler(msg.Data)
		msg.Ack() // Acknowledgment para confiabilidade
	})
	return err
}

func (q *NatsQueue) SubscribeQueueWithWorkers(topic, queueGroup string, handler func(data []byte), concurrency int) error {
	msgChan := make(chan []byte, 10000) // Aumentar buffer para 10k mensagens

	// Inicia os workers
	for i := 0; i < concurrency; i++ {
		go func() {
			for data := range msgChan {
				handler(data)
			}
		}()
	}

	// NATS assina e envia para o canal
	_, err := q.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {
		select {
		case msgChan <- msg.Data:
			// Mensagem enviada com sucesso
		default:
			// Canal cheio, descartar mensagem para evitar bloqueio
			log.Printf("Warning: message channel full, dropping message")
		}
		msg.Ack() // Acknowledgment
	})

	return err
}
