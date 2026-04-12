package port

import "context"

// MessageHandler is the callback invoked for each message received from the queue.
type MessageHandler func(ctx context.Context, ev OrderEvent) error

// MessageConsumer reads from a queue and dispatches to a MessageHandler.
type MessageConsumer interface {
	Consume(ctx context.Context, queue string, handler MessageHandler) error
	Close() error
}
