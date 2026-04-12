package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/renan/go-api-worker/internal/application/port"
)

// EnqueueEmail handles a Kafka OrderEvent and forwards it to a RabbitMQ queue.
// It implements port.EventHandler.
type EnqueueEmail struct {
	Publisher port.MessagePublisher
	Queue     string
}

func (uc EnqueueEmail) HandleEvent(ctx context.Context, ev port.OrderEvent) error {
	log.Printf("enqueue email for order %s status %s", ev.OrderID, ev.Status)
	if err := uc.Publisher.PublishMessage(ctx, uc.Queue, ev); err != nil {
		return fmt.Errorf("enqueue email: %w", err)
	}
	return nil
}
