package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renan/go-api-worker/internal/application/usecase"
	"github.com/renan/go-api-worker/internal/infra/config"
	kafkainfra "github.com/renan/go-api-worker/internal/infra/kafka"
	"github.com/renan/go-api-worker/internal/infra/rabbitmq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Printf("email-worker error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadEmailWorker()
	if err != nil {
		return err
	}

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	defer publisher.Close()

	rmqConsumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	defer rmqConsumer.Close()

	enqueueUC := usecase.EnqueueEmail{
		Publisher: publisher,
		Queue:     cfg.RabbitMQQueue,
	}
	sendUC := usecase.SendEmail{
		Sleep: time.Sleep,
		Delay: 2 * time.Second,
	}

	kafkaConsumer := kafkainfra.NewConsumer(cfg.KafkaBrokers)

	errc := make(chan error, 2)

	go func() {
		log.Printf("email-worker: kafka consumer started topic=%s group=email-worker", cfg.KafkaTopic)
		errc <- kafkaConsumer.Consume(ctx, cfg.KafkaTopic, "email-worker", enqueueUC.HandleEvent)
	}()

	go func() {
		log.Printf("email-worker: rabbitmq consumer started queue=%s", cfg.RabbitMQQueue)
		errc <- rmqConsumer.Consume(ctx, cfg.RabbitMQQueue, sendUC.HandleMessage)
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return nil
	}
}
