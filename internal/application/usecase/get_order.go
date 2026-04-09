package usecase

import (
	"context"
	"fmt"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
)

type GetOrder struct {
	Repo port.OrderRepository
}

func (uc GetOrder) Execute(ctx context.Context, orderID string) (order.Order, error) {
	if uc.Repo == nil {
		return order.Order{}, fmt.Errorf("repo is required")
	}
	if orderID == "" {
		return order.Order{}, fmt.Errorf("order_id is required")
	}
	return uc.Repo.GetByID(ctx, orderID)
}

