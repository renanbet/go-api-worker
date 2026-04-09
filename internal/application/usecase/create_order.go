package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
)

type CreateOrder struct {
	Repo      port.OrderRepository
	Publisher port.EventPublisher
	Topic     string
	Now       func() time.Time
}

type CreateOrderResult struct {
	OrderID string
	Status  order.Status
}

func (uc CreateOrder) Execute(ctx context.Context, product string, quantity int) (CreateOrderResult, error) {
	if stringsTrim(product) == "" {
		return CreateOrderResult{}, fmt.Errorf("product is required")
	}
	if quantity <= 0 {
		return CreateOrderResult{}, fmt.Errorf("quantity must be > 0")
	}
	if uc.Repo == nil {
		return CreateOrderResult{}, fmt.Errorf("repo is required")
	}
	if uc.Publisher == nil {
		return CreateOrderResult{}, fmt.Errorf("publisher is required")
	}
	if uc.Topic == "" {
		return CreateOrderResult{}, fmt.Errorf("topic is required")
	}
	now := uc.Now
	if now == nil {
		now = time.Now
	}

	orderID := uuid.NewString()
	o := order.Order{
		OrderID:     orderID,
		ProductName: product,
		Quantity:    quantity,
		Status:      order.StatusCreated,
		CreatedAt:   now(),
	}
	if err := uc.Repo.Create(ctx, o); err != nil {
		return CreateOrderResult{}, err
	}

	ev := port.OrderEvent{OrderID: orderID, Status: string(order.StatusProcessing)}
	if err := uc.Publisher.PublishOrderEvent(ctx, uc.Topic, ev); err != nil {
		return CreateOrderResult{}, err
	}

	return CreateOrderResult{OrderID: orderID, Status: o.Status}, nil
}

func stringsTrim(s string) string {
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\n' || s[i] == '\t' || s[i] == '\r') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\n' || s[j-1] == '\t' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}

