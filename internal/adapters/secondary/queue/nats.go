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
		msg.Ack()
	})
	return err
}

func (q *NatsQueue) SubscribeQueueWithWorkers(topic, queueGroup string, handler func(data []byte), concurrency int) error {
	msgChan := make(chan []byte, 10000)

	for i := 0; i < concurrency; i++ {
		go func() {
			for data := range msgChan {
				handler(data)
			}
		}()
	}

	_, err := q.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {
		select {
		case msgChan <- msg.Data:
		default:
			log.Printf("Warning: message channel full, dropping message")
		}
		msg.Ack()
	})

	return err
}
