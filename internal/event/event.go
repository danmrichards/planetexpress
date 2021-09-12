package event

import (
	"fmt"
	"strconv"
)

const (
	// PackageAllocate occurs when space for a package is allocated on the ship.
	PackageAllocate PackageEventType = "package_allocate"

	// PackageLoad occurs when a package is loaded onto the ship.
	PackageLoad PackageEventType = "package_load"

	// PackageUnload occurs when a package is unloaded from the ship.
	PackageUnload PackageEventType = "package_unload"
)

// PackageEventType represents the type of event.
type PackageEventType string

// PackageEvent represents an event occurring on a single package.
type PackageEvent struct {
	EventID     string
	Typ         PackageEventType
	PackageID   string
	PackageSize int
}

// Values returns a map representation of the event for use with an event
// dispatcher.
func (p PackageEvent) Values() map[string]interface{} {
	return map[string]interface{}{
		"package_id":   p.PackageID,
		"type":         string(p.Typ),
		"package_size": p.PackageSize,
	}
}

// FromValues returns a new PackageEvent from the given values.
func FromValues(id string, values map[string]interface{}) (*PackageEvent, error) {
	typVal, ok := values["type"]
	if !ok {
		return nil, MissingValueError("type")
	}
	typ, ok := typVal.(string)
	if !ok {
		return nil, NewValueTypeError("type", typVal)
	}

	pidVal, ok := values["package_id"]
	if !ok {
		return nil, MissingValueError("package_id")
	}
	pid, ok := pidVal.(string)
	if !ok {
		return nil, NewValueTypeError("package_id", pidVal)
	}

	sizeVal, ok := values["package_size"]
	if !ok {
		return nil, MissingValueError("package_size")
	}
	sizeStr, ok := sizeVal.(string)
	if !ok {
		return nil, NewValueTypeError("package_size", sizeVal)
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("parse size: %w", err)
	}

	return &PackageEvent{
		EventID:     id,
		Typ:         PackageEventType(typ),
		PackageID:   pid,
		PackageSize: size,
	}, nil
}
