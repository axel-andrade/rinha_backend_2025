package queue

import (
	"log"
	"os"
	"sync"

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

func (q *NatsQueue) Publish(topic string, data []byte) error {
	return q.conn.Publish(topic, data)
}

func (q *NatsQueue) SubscribeQueue(topic, queueGroup string, handler func(data []byte)) error {
	_, err := q.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}

func (q *NatsQueue) SubscribeQueueWithWorkers(topic, queueGroup string, handler func(data []byte), concurrency int) error {
	msgChan := make(chan []byte, 100) // buffer de mensagens

	// inicia os workers
	for i := 0; i < concurrency; i++ {
		go func() {
			for data := range msgChan {
				handler(data)
			}
		}()
	}

	// NATS assina e envia para o canal
	_, err := q.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {
		msgChan <- msg.Data
	})

	return err
}
