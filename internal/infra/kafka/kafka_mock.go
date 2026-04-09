package kafka

import (
	"context"
	"github.com/renan/go-api-worker/internal/application/port"
)

type EventPublisherMock struct {
	PublishOrderEventFunc func(ctx context.Context, topic string, event port.OrderEvent) error
	PublishOrderEventCount int
}

func (m *EventPublisherMock) PublishOrderEvent(ctx context.Context, topic string, event port.OrderEvent) error {
	m.PublishOrderEventCount++
	return m.PublishOrderEventFunc(ctx, topic, event)
}



