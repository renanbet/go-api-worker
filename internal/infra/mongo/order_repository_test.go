package mongo

import (
	"testing"

	"github.com/renan/go-api-worker/internal/application/port"
)

func TestOrderRepository_ImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ port.OrderRepository = (*orderRepository)(nil)
}

