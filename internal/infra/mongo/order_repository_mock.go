package mongo

import (
	"context"

	"github.com/renan/go-api-worker/internal/domain/order"
)

type OrderRepositoryMock struct {
	CreateFunc        func(ctx context.Context, o order.Order) error
	CreateCount       int
	UpdateStatusFunc  func(ctx context.Context, orderID string, status order.Status) error
	UpdateStatusCount int
	GetByIDFunc       func(ctx context.Context, orderID string) (order.Order, error)
	GetByIDCount      int
}

func (m *OrderRepositoryMock) Create(ctx context.Context, o order.Order) error {
	m.CreateCount++
	return m.CreateFunc(ctx, o)
}

func (m *OrderRepositoryMock) UpdateStatus(ctx context.Context, orderID string, status order.Status) error {
	m.UpdateStatusCount++
	return m.UpdateStatusFunc(ctx, orderID, status)
}

func (m *OrderRepositoryMock) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	m.GetByIDCount++
	return m.GetByIDFunc(ctx, orderID)
}
