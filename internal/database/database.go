package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// modernc.org/sqlite is a pure-Go translation of SQLite, chosen over
	// mattn/go-sqlite3 because mattn requires cgo. With cgo enabled the build
	// produces a dynamically linked binary that has to match the runtime
	// image's libc, which forces a fat runtime image and breaks the
	// single-static-binary deployment this repo is built around. The pure-Go
	// driver trades some raw speed for CGO_ENABLED=0; at this site's traffic
	// the speed is irrelevant and the simpler build is not.
	//
	// Note the driver name is "sqlite", not "sqlite3".
	_ "modernc.org/sqlite"
)

// Service is the surface the rest of the app depends on. Kept narrow so the
// storage layer can be swapped (or dropped into a different site in the
// multi-site plan) without the server package knowing.
type Service interface {
	Health(ctx context.Context) error
	RecordPageview(ctx context.Context, path string, status int, dur time.Duration, referrer string) error
	PageviewSummary(ctx context.Context, days int) (Summary, error)
	Close() error
}

type service struct {
	db *sql.DB
}

// New opens the database at dburl and verifies it is reachable.
//
// This replaces a package-level singleton that returned a cached *service and
// swallowed re-open attempts. The singleton made the failure mode worse rather
// than better: it hid configuration mistakes behind a stale handle, and its
// unguarded read of a package variable raced if it were ever called from more
// than one goroutine. Callers construct one instance at startup and pass it
// down, which is both simpler and testable — a test can open its own :memory:
// database without stomping on a global.
func New(dburl string) (Service, error) {
	db, err := sql.Open("sqlite", dburl)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", dburl, err)
	}

	// SQLite's PRAGMA foreign_keys is per-connection and database/sql pools
	// connections transparently, so with a pool of N a cascade delete would
	// apply only on whichever connection happened to have run the pragma.
	// SQLite serializes writes regardless, so at this scale the pool buys
	// nothing anyway.
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite %q: %w", dburl, err)
	}

	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &service{db: db}, nil
}

// Health reports whether the database is reachable.
//
// It returns an error rather than calling log.Fatalf, as the previous version
// did. A health check that kills the process on a transient failure turns a
// blip that a caller could have retried into an outage of the entire site —
// including the pages that do not touch the database at all.
func (s *service) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database unreachable: %w", err)
	}
	return nil
}

func (s *service) Close() error {
	return s.db.Close()
}
