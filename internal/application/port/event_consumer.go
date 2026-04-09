package port

import "context"

type EventHandler func(ctx context.Context, event OrderEvent) error

type EventConsumer interface {
	Consume(ctx context.Context, topic string, groupID string, handler EventHandler) error
}

