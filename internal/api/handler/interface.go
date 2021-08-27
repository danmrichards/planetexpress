package handler

import "context"

// EventService is the interface implemented by types that can act as a service
// for managing events.
type EventService interface {
	// PackageAllocate dispatches an event indicating the allocation of a given
	// package in the ship cargo bay.
	PackageAllocate(ctx context.Context, id string, size int) error
}

// pkgIDFunc is a function that returns a generated package ID.
type pkgIDFunc func() string
