package port

import "testing"

func TestErrNotFoundImplementsError(t *testing.T) {
	t.Parallel()

	var _ error = (ErrNotFound{})
}
