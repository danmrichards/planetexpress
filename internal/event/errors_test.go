package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissingValueError_Error(t *testing.T) {
	require.Equal(
		t,
		MissingValueError("foo").Error(),
		`missing value "foo"`,
	)
}

func TestValueTypeError_Error(t *testing.T) {
	require.Equal(
		t,
		ValueTypeError{"foo", 123}.Error(),
		`unexpected type (int) for value "foo"`,
	)
}
