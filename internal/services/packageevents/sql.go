package packageevents

import (
	"database/sql"
	"fmt"

	"github.com/danmrichards/planetexpress/internal/event"
)

// SQLService is a service for returning package events, backed by SQL.
type SQLService struct {
	db *sql.DB
}

// NewSQLService returns a new service using the given DB handle.
func NewSQLService(db *sql.DB) *SQLService {
	return &SQLService{db}
}

// Latest returns the most recent package event.
func (s *SQLService) Latest() (*event.PackageEvent, error) {
	pkgEvt := &event.PackageEvent{}
	if err := s.db.QueryRow(
		`SELECT event_id, event_type, package_id, package_size
FROM package_events
ORDER BY event_id DESC`).Scan(
		&pkgEvt.EventID, &pkgEvt.Typ, &pkgEvt.PackageID, &pkgEvt.PackageSize); err != nil {
		return nil, fmt.Errorf("get latest package event: %w", err)
	}

	return pkgEvt, nil
}

// Save saves a package event.
func (s *SQLService) Save(evt *event.PackageEvent) error {
	_, err := s.db.Exec(
		`INSERT INTO package_events (event_id, event_type, package_id, package_size)
VALUES ($1, $2, $3, $4)`,
		evt.EventID,
		evt.Typ,
		evt.PackageID,
		evt.PackageSize,
	)
	if err != nil {
		return fmt.Errorf("save event: %w", err)
	}

	return nil
}

// UpdateShipStatus updates the status of the ship based on the received event.
func (s *SQLService) UpdateShipStatus(evt *event.PackageEvent) (err error) {
	switch evt.Typ {
	case event.PackageAllocate:
		return s.allocatePackage(evt.PackageSize)
	case event.PackageLoad:
		return s.loadPackage(evt.PackageSize)
	case event.PackageUnload:
		return s.unloadPackage(evt.PackageSize)
	default:
		return UnknownEventError(evt.Typ)
	}
}

func (s *SQLService) allocatePackage(size int) error {
	_, err := s.db.Exec(
		"UPDATE ship_status SET allocated = allocated + $1, available = available - $1",
		size,
	)
	if err != nil {
		return fmt.Errorf("allocate package: %w", err)
	}

	return nil
}

func (s *SQLService) loadPackage(size int) error {
	_, err := s.db.Exec(
		"UPDATE ship_status SET loaded = loaded + $1, allocated = allocated - $1",
		size,
	)
	if err != nil {
		return fmt.Errorf("load package: %w", err)
	}

	return nil
}

func (s *SQLService) unloadPackage(size int) error {
	_, err := s.db.Exec(
		"UPDATE ship_status SET loaded = loaded - $1, available = available + $1",
		size,
	)
	if err != nil {
		return fmt.Errorf("unload package: %w", err)
	}

	return nil
}
