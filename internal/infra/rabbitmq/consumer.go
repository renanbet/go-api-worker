package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/renan/go-api-worker/internal/application/port"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq consumer dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq consumer channel: %w", err)
	}
	if err := ch.Qos(1, 0, false); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq qos: %w", err)
	}
	return &Consumer{conn: conn, channel: ch}, nil
}

func (c *Consumer) Consume(ctx context.Context, queue string, handler port.MessageHandler) error {
	_, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("channel consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq channel closed")
			}
			var ev port.OrderEvent
			if err := json.Unmarshal(d.Body, &ev); err != nil {
				log.Printf("rabbitmq: decode error, discarding message: %v", err)
				_ = d.Nack(false, false)
				continue
			}
			if err := handler(ctx, ev); err != nil {
				log.Printf("rabbitmq: handler error, requeuing: %v", err)
				_ = d.Nack(false, true)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			_ = d.Ack(false)
		}
	}
}

func (c *Consumer) Close() error {
	_ = c.channel.Close()
	return c.conn.Close()
}
