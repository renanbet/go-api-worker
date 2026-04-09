package kafka

import (
	"context"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		RequiredAcks: kafka.RequireOne,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
	}
	return &Producer{writer: w}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) PublishOrderEvent(ctx context.Context, topic string, event port.OrderEvent) error {
	value, err := encodeOrderEvent(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.OrderID),
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, msg)
}

