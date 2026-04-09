package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
	"github.com/stretchr/testify/assert"
)

type getRepoStub struct {
	o   order.Order
	err error
}

func (g getRepoStub) Create(ctx context.Context, o order.Order) error { return nil }
func (g getRepoStub) UpdateStatus(ctx context.Context, orderID string, status order.Status) error {
	return nil
}
func (g getRepoStub) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	return g.o, g.err
}

func TestGetOrder_ReturnsOrder(t *testing.T) {
	t.Parallel()

	want := order.Order{
		OrderID:     "id1",
		ProductName: "p",
		Quantity:    1,
		Status:      order.StatusCreated,
		CreatedAt:   time.Unix(0, 0).UTC(),
	}

	uc := GetOrder{Repo: getRepoStub{o: want}}
	got, err := uc.Execute(context.Background(), "id1")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetOrder_ValidatesInput(t *testing.T) {
	t.Parallel()

	uc := GetOrder{Repo: getRepoStub{}}
	_, err := uc.Execute(context.Background(), "")
	assert.Error(t, err)

	uc2 := GetOrder{}
	_, err = uc2.Execute(context.Background(), "x")
	assert.Error(t, err)
}

func TestGetOrder_NotFoundPassThrough(t *testing.T) {
	t.Parallel()

	uc := GetOrder{Repo: getRepoStub{err: port.ErrNotFound{}}}
	_, err := uc.Execute(context.Background(), "x")
	if err == nil {
		t.Fatalf("expected error")
	}
}
