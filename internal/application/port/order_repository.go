package port

import (
	"context"

	"github.com/renan/go-api-worker/internal/domain/order"
)

type OrderRepository interface {
	Create(ctx context.Context, o order.Order) error
	GetByID(ctx context.Context, orderID string) (order.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status order.Status) error
}

