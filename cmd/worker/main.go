package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/renan/go-api-worker/internal/application/usecase"
	"github.com/renan/go-api-worker/internal/infra/config"
	kafkainfra "github.com/renan/go-api-worker/internal/infra/kafka"
	mongoinfra "github.com/renan/go-api-worker/internal/infra/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Printf("worker error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	mc, err := mongoConnect(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer func() { _ = mc.Disconnect(context.Background()) }()

	repo, err := mongoinfra.NewOrderRepository(mc.Database(cfg.MongoDB))
	if err != nil {
		return err
	}

	consumer := kafkainfra.NewConsumer(cfg.KafkaBrokers)
	uc := usecase.ProcessOrder{Repo: repo}

	log.Printf("worker consuming topic=%s group=%s", cfg.KafkaTopic, cfg.KafkaGroupID)
	return consumer.Consume(ctx, cfg.KafkaTopic, cfg.KafkaGroupID, uc.HandleEvent)
}

func mongoConnect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return client, nil
}

