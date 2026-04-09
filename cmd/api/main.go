package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renan/go-api-worker/internal/application/usecase"
	httpinfra "github.com/renan/go-api-worker/internal/infra/http"
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
		log.Printf("api error: %v", err)
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

	producer := kafkainfra.NewProducer(cfg.KafkaBrokers)
	defer func() { _ = producer.Close() }()

	createUC := usecase.CreateOrder{
		Repo:      repo,
		Publisher: producer,
		Topic:     cfg.KafkaTopic,
	}

	h := &httpinfra.Handler{CreateOrderUC: createUC}
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpinfra.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("api listening on %s", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func mongoConnect(ctx context.Context, uri string) (*mongo.Client, error) {
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(cctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(cctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return client, nil
}

