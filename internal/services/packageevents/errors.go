package packageevents

import "fmt"

// UnknownEventError is returned when an event is found with an unknown type.
type UnknownEventError string

func (u UnknownEventError) Error() string {
	return fmt.Sprintf("unknown event type %q", string(u))
}
