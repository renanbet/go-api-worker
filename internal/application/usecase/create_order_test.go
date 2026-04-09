package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
)

type mockRepo struct {
	created []order.Order
	updates []struct {
		id     string
		status order.Status
	}
}

func (m *mockRepo) Create(ctx context.Context, o order.Order) error {
	m.created = append(m.created, o)
	return nil
}

func (m *mockRepo) UpdateStatus(ctx context.Context, orderID string, status order.Status) error {
	m.updates = append(m.updates, struct {
		id     string
		status order.Status
	}{orderID, status})
	return nil
}

type mockPublisher struct {
	events []port.OrderEvent
	topics []string
}

func (m *mockPublisher) PublishOrderEvent(ctx context.Context, topic string, event port.OrderEvent) error {
	m.topics = append(m.topics, topic)
	m.events = append(m.events, event)
	return nil
}

func TestCreateOrder_StoresAndPublishes(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{}
	pub := &mockPublisher{}
	now := time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC)

	uc := CreateOrder{
		Repo:      repo,
		Publisher: pub,
		Topic:     "order_events",
		Now:       func() time.Time { return now },
	}

	res, err := uc.Execute(context.Background(), "mouse", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrderID == "" {
		t.Fatalf("expected order id")
	}
	if res.Status != order.StatusCreated {
		t.Fatalf("expected CREATED, got %s", res.Status)
	}

	if len(repo.created) != 1 {
		t.Fatalf("expected 1 created order, got %d", len(repo.created))
	}
	if repo.created[0].OrderID != res.OrderID {
		t.Fatalf("stored orderID mismatch")
	}
	if repo.created[0].CreatedAt != now {
		t.Fatalf("expected created_at to be set")
	}
	if repo.created[0].Status != order.StatusCreated {
		t.Fatalf("expected stored CREATED")
	}

	if len(pub.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(pub.events))
	}
	if pub.events[0].OrderID != res.OrderID {
		t.Fatalf("event order id mismatch")
	}
	if pub.events[0].Status != string(order.StatusProcessing) {
		t.Fatalf("expected PROCESSING in event, got %s", pub.events[0].Status)
	}
}

func TestCreateOrder_ValidatesInput(t *testing.T) {
	t.Parallel()

	uc := CreateOrder{
		Repo:      &mockRepo{},
		Publisher: &mockPublisher{},
		Topic:     "order_events",
	}

	if _, err := uc.Execute(context.Background(), "", 1); err == nil {
		t.Fatalf("expected error for empty product")
	}
	if _, err := uc.Execute(context.Background(), "x", 0); err == nil {
		t.Fatalf("expected error for quantity 0")
	}
}

