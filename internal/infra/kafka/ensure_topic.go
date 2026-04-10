package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// EnsureTopic cria o tópico no cluster se ainda não existir (idempotente).
// Usa o broker controller (CreateTopics); adequado ao compose com 1 broker (RF=1).
func EnsureTopic(ctx context.Context, brokers []string, topic string, numPartitions int) error {
	if len(brokers) == 0 {
		return fmt.Errorf("kafka: no brokers")
	}
	if topic == "" {
		return fmt.Errorf("kafka: topic is required")
	}
	if numPartitions <= 0 {
		return fmt.Errorf("kafka: numPartitions must be > 0")
	}

	broker := brokers[0]
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		return fmt.Errorf("kafka dial %q: %w", broker, err)
	}

	controller, err := conn.Controller()
	_ = conn.Close()
	if err != nil {
		return fmt.Errorf("kafka controller: %w", err)
	}

	addr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	ctrlConn, err := kafka.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("kafka dial controller %q: %w", addr, err)
	}
	defer func() { _ = ctrlConn.Close() }()

	err = ctrlConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("kafka create topic %q: %w", topic, err)
	}
	return nil
}
