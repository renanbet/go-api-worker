package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/renan/go-api-worker/internal/application/port"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq publisher dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq publisher channel: %w", err)
	}
	return &Publisher{conn: conn, channel: ch}, nil
}

func (p *Publisher) PublishMessage(ctx context.Context, queue string, ev port.OrderEvent) error {
	_, err := p.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}
	body, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	return p.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (p *Publisher) Close() error {
	_ = p.channel.Close()
	return p.conn.Close()
}
