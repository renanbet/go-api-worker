package kafka

import (
	"testing"

	"github.com/renan/go-api-worker/internal/application/port"
)

func TestEncodeDecodeOrderEvent_RoundTrip(t *testing.T) {
	t.Parallel()

	in := port.OrderEvent{OrderID: "id1", Status: "PROCESSING"}
	b, err := encodeOrderEvent(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := decodeOrderEvent(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != in {
		t.Fatalf("round trip mismatch: %+v != %+v", out, in)
	}
}

