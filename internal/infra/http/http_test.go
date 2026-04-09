package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/application/usecase"
	"github.com/renan/go-api-worker/internal/domain/order"
)

type stubRepo struct{}

func (s stubRepo) Create(ctx context.Context, o order.Order) error { return nil }
func (s stubRepo) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	return order.Order{OrderID: orderID}, nil
}
func (s stubRepo) UpdateStatus(ctx context.Context, orderID string, status order.Status) error {
	return nil
}

type stubPublisher struct{}

func (s stubPublisher) PublishOrderEvent(ctx context.Context, topic string, event port.OrderEvent) error {
	return nil
}

func TestCreateOrderHandler_Returns201(t *testing.T) {
	t.Parallel()

	h := &Handler{
		CreateOrderUC: usecase.CreateOrder{
			Repo:      stubRepo{},
			Publisher: stubPublisher{},
			Topic:     "order_events",
			Now:       func() time.Time { return time.Unix(0, 0).UTC() },
		},
		GetOrderUC: usecase.GetOrder{Repo: stubRepo{}},
	}

	srv := httptest.NewServer(NewRouter(h))
	t.Cleanup(srv.Close)

	resp, err := http.Post(srv.URL+"/orders", "application/json", bytes.NewBufferString(`{"product":"p","quantity":1}`))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

type notFoundRepo struct{ stubRepo }

func (n notFoundRepo) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	return order.Order{}, port.ErrNotFound{}
}

func TestGetOrderHandler_Returns200(t *testing.T) {
	t.Parallel()

	h := &Handler{
		CreateOrderUC: usecase.CreateOrder{
			Repo:      stubRepo{},
			Publisher: stubPublisher{},
			Topic:     "order_events",
			Now:       func() time.Time { return time.Unix(0, 0).UTC() },
		},
		GetOrderUC: usecase.GetOrder{Repo: stubRepo{}},
	}

	srv := httptest.NewServer(NewRouter(h))
	t.Cleanup(srv.Close)

	resp, err := http.Get(srv.URL + "/orders/abc")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetOrderHandler_Returns404(t *testing.T) {
	t.Parallel()

	h := &Handler{
		CreateOrderUC: usecase.CreateOrder{
			Repo:      stubRepo{},
			Publisher: stubPublisher{},
			Topic:     "order_events",
			Now:       func() time.Time { return time.Unix(0, 0).UTC() },
		},
		GetOrderUC: usecase.GetOrder{Repo: notFoundRepo{}},
	}

	srv := httptest.NewServer(NewRouter(h))
	t.Cleanup(srv.Close)

	resp, err := http.Get(srv.URL + "/orders/missing")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

