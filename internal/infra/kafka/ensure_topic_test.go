package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureTopic_ValidatesInput(t *testing.T) {
	t.Parallel()

	err := EnsureTopic(context.Background(), nil, "x")
	assert.Error(t, err)
	err = EnsureTopic(context.Background(), []string{"localhost:9092"}, "")
	assert.Error(t, err)
}
