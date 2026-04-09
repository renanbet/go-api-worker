package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
	"github.com/renan/go-api-worker/internal/infra/mongo"
	"github.com/stretchr/testify/assert"
)

func TestProcessOrder_UpdatesProcessingAndConcluded(t *testing.T) {
	t.Parallel()

	repo := &mongo.OrderRepositoryMock{
		UpdateStatusFunc: func(ctx context.Context, orderID string, status order.Status) error {
			return nil
		},
	}
	var slept time.Duration

	uc := ProcessOrder{
		Repo:  repo,
		Sleep: func(d time.Duration) { slept = d },
		Delay: 2 * time.Second,
	}

	err := uc.HandleEvent(context.Background(), port.OrderEvent{OrderID: "abc", Status: "PROCESSING"})
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, slept)
	assert.Equal(t, 2, repo.UpdateStatusCount)
}
