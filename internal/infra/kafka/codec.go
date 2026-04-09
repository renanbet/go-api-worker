package kafka

import (
	"encoding/json"

	"github.com/renan/go-api-worker/internal/application/port"
)

func encodeOrderEvent(ev port.OrderEvent) ([]byte, error) {
	return json.Marshal(ev)
}

func decodeOrderEvent(b []byte) (port.OrderEvent, error) {
	var ev port.OrderEvent
	err := json.Unmarshal(b, &ev)
	return ev, err
}

