package port

import "context"

type OrderEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type EventPublisher interface {
	PublishOrderEvent(ctx context.Context, topic string, event OrderEvent) error
}

