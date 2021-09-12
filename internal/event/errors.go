package event

import "fmt"

// MissingValueError is returned when an expected event value is missing.
type MissingValueError string

func (m MissingValueError) Error() string {
	return fmt.Sprintf("missing value: %q", string(m))
}

// ValueTypeError is returned when an expected event value is of the wrong type.
type ValueTypeError struct {
	val string
	v   interface{}
}

// NewValueTypeError returns a new ValueTypeError.
func NewValueTypeError(val string, v interface{}) error {
	return ValueTypeError{
		val: val,
		v:   v,
	}
}

func (m ValueTypeError) Error() string {
	return fmt.Sprintf("unexpected type (%[1]T) for value %q", m.v, m.val)
}
