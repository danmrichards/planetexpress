package event

// PackageEventType represents the type of event.
type PackageEventType string

// PackageEvent represents an event occurring on a single package.
type PackageEvent struct {
	ID   string
	Typ  PackageEventType
	Size int
}

// Values returns a map representation of the event for use with an event
// dispatcher.
func (p PackageEvent) Values() map[string]interface{} {
	return map[string]interface{}{
		"package_id": p.ID,
		"type":       string(p.Typ),
		"size":       p.Size,
	}
}

const (
	// PackageAllocate occurs when space for a package is allocated on the ship.
	PackageAllocate PackageEventType = "package_allocate"

	// PackageLoad occurs when a package is loaded onto the ship.
	PackageLoad PackageEventType = "package_load"

	// PackageUnload occurs when a package is unloaded from the ship.
	PackageUnload PackageEventType = "package_unload"
)
