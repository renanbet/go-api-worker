package usecase

import (
	"context"
	"log"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
)

// SendEmail simulates sending an email for a given order event.
// It implements the port.MessageHandler function signature.
type SendEmail struct {
	Sleep Sleeper      // reuses Sleeper type from process_order.go
	Delay time.Duration
}

func (uc SendEmail) HandleMessage(ctx context.Context, ev port.OrderEvent) error {
	sleep := uc.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	delay := uc.Delay
	if delay <= 0 {
		delay = 2 * time.Second
	}
	log.Printf("sending email for order %s status %s", ev.OrderID, ev.Status)
	sleep(delay)
	log.Printf("email sent for order %s", ev.OrderID)
	return nil
}
