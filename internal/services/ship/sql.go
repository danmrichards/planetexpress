package ship

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/danmrichards/planetexpress/internal/event"
)

type shipAvailability struct {
	Available int
	UpdatedAt time.Time
}

// SQLService is a service for returning ship status, backed by SQL.
type SQLService struct {
	db *sql.DB
}

// NewSQLService returns a new service using the given DB handle.
func NewSQLService(db *sql.DB) *SQLService {
	return &SQLService{db}
}

// CapacityAvailable implements handler.ShipService
func (s *SQLService) CapacityAvailable(size int) (bool, error) {
	sa := &shipAvailability{}
	if err := s.db.QueryRow(
		"SELECT available, updated_at FROM ship_status",
	).Scan(&sa.Available, &sa.UpdatedAt); err != nil {
		return false, fmt.Errorf("ship availability: %w", err)
	}

	// Ensure we include capacity updates from any events not yet persisted
	// into the ship_status view.
	events, err := s.eventsAfter(sa.UpdatedAt)
	if err != nil {
		return false, fmt.Errorf("could not get events: %w", err)
	}

	// Adjust ship status from the new events.
	for _, evt := range events {
		switch evt.Typ {
		case event.PackageAllocate:
			sa.Available -= evt.PackageSize
		case event.PackageUnload:
			sa.Available += evt.PackageSize
		}
	}

	return sa.Available >= size, nil
}

// eventsAfter returns the list of package events created after a given time.
func (s *SQLService) eventsAfter(at time.Time) ([]event.PackageEvent, error) {
	rows, err := s.db.Query(
		`SELECT event_id, event_type, package_id, package_size
FROM package_events 
WHERE created_at > $1
ORDER BY created_at ASC`,
		at,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.PackageEvent
	for rows.Next() {
		var evt event.PackageEvent
		if err = rows.Scan(
			&evt.EventID, &evt.Typ, &evt.PackageID, &evt.PackageSize,
		); err != nil {
			return nil, err
		}
		events = append(events, evt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
