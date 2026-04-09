package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	brokers []string
}

func NewConsumer(brokers []string) *Consumer {
	return &Consumer{brokers: brokers}
}

func (c *Consumer) Consume(ctx context.Context, topic string, groupID string, handler port.EventHandler) error {
	if len(c.brokers) == 0 {
		return fmt.Errorf("brokers are required")
	}
	if topic == "" {
		return fmt.Errorf("topic is required")
	}
	if groupID == "" {
		return fmt.Errorf("groupID is required")
	}
	if handler == nil {
		return fmt.Errorf("handler is required")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0, // manual commit
	})
	defer r.Close()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			return err
		}
		ev, err := decodeOrderEvent(m.Value)
		if err == nil {
			err = handler(ctx, ev)
		}
		if err == nil {
			err = r.CommitMessages(ctx, m)
		}
		if err != nil {
			_ = r.SetOffset(m.Offset) // best-effort: allow retry
			time.Sleep(250 * time.Millisecond)
		}
	}
}

