package packageevents

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnknownEventError_Error(t *testing.T) {
	require.Equal(
		t,
		UnknownEventError("foo").Error(),
		`unknown event type "foo"`,
	)
}
