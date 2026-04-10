package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
)

type Sleeper func(time.Duration)

type ProcessOrder struct {
	Repo  port.OrderRepository
	Sleep Sleeper
	Delay time.Duration
}

func (uc ProcessOrder) HandleEvent(ctx context.Context, ev port.OrderEvent) error {
	log.Printf("processing order %s with status %s", ev.OrderID, ev.Status)
	if uc.Repo == nil {
		return fmt.Errorf("repo is required")
	}
	if ev.OrderID == "" {
		return fmt.Errorf("order_id is required")
	}

	sleep := uc.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	delay := uc.Delay
	if delay <= 0 {
		delay = 2 * time.Second
	}

	if err := uc.Repo.UpdateStatus(ctx, ev.OrderID, order.StatusProcessing); err != nil {
		return err
	}
	sleep(delay)
	if err := uc.Repo.UpdateStatus(ctx, ev.OrderID, order.StatusConcluded); err != nil {
		return err
	}
	log.Printf("order %s processed with status %s", ev.OrderID, order.StatusConcluded)
	return nil
}
