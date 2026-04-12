package port

import "context"

// MessagePublisher publishes an OrderEvent to a message queue (e.g. RabbitMQ).
type MessagePublisher interface {
	PublishMessage(ctx context.Context, queue string, ev OrderEvent) error
	Close() error
}
