package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
)

func TestProcessOrder_UpdatesProcessingAndConcluded(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{}
	var slept time.Duration

	uc := ProcessOrder{
		Repo:  repo,
		Sleep: func(d time.Duration) { slept = d },
		Delay: 2 * time.Second,
	}

	err := uc.HandleEvent(context.Background(), port.OrderEvent{OrderID: "abc", Status: "PROCESSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if slept != 2*time.Second {
		t.Fatalf("expected sleep 2s, got %v", slept)
	}
	if len(repo.updates) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(repo.updates))
	}
	if repo.updates[0].status != order.StatusProcessing {
		t.Fatalf("expected PROCESSING, got %s", repo.updates[0].status)
	}
	if repo.updates[1].status != order.StatusConcluded {
		t.Fatalf("expected CONCLUDED, got %s", repo.updates[1].status)
	}
}

