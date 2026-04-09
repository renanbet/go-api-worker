package order

import "testing"

func TestStatusValues(t *testing.T) {
	t.Run("values must be non-empty", func(t *testing.T) {
		if StatusCreated == "" || StatusProcessing == "" || StatusConcluded == "" {
			t.Fatalf("status values must be non-empty")
		}
	})
}

func TestOrderValidate(t *testing.T) {
	t.Run("valid order", func(t *testing.T) {
		o := Order{ProductName: "Product", Quantity: 1, Status: StatusCreated}
		if err := o.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
	t.Run("invalid order", func(t *testing.T) {
		o := Order{ProductName: "", Quantity: 0, Status: StatusCreated}
		if err := o.Validate(); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
