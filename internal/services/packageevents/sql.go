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
