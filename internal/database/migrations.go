package database

import (
	"context"
	"database/sql"
	"fmt"
)

// migrations are applied in order, exactly once each, tracked by index in the
// schema_migrations table. Append only — editing a migration that has already
// run on the droplet will not re-run it, so the two databases would silently
// diverge.
var migrations = []string{
	// duration_us, not duration_ms. These pages are rendered from templates in
	// memory and complete in well under a millisecond, so whole milliseconds
	// rounded every single request to 0 and the percentile column was
	// structurally incapable of reporting anything. Microseconds give three
	// more digits of headroom; the API still reports milliseconds, as a float.
	`CREATE TABLE IF NOT EXISTS pageviews (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		path        TEXT    NOT NULL,
		status      INTEGER NOT NULL,
		duration_us INTEGER NOT NULL,
		referrer    TEXT    NOT NULL DEFAULT '',
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE INDEX IF NOT EXISTS idx_pageviews_created_at ON pageviews(created_at)`,
	`CREATE INDEX IF NOT EXISTS idx_pageviews_path ON pageviews(path)`,
}

// migrate brings the schema up to date. It runs at startup rather than as a
// separate command because this is a single-binary deployment with one writer
// and a schema measured in one table — a migration tool would be more moving
// parts than the thing it manages.
//
// Deliberately absent from the schema: IP addresses and user agents. This is a
// personal site; there is no question they would answer that justifies holding
// visitor identifiers, and not collecting them is a stronger claim than any
// retention policy. Referrer is kept because it is the one field that answers
// the question the metrics exist for — whether a link shared into Slack or
// LinkedIn actually brought anyone — and it identifies the source, not the
// person.
func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	var applied int
	// COALESCE because MAX over zero rows is NULL, which will not scan into an
	// int. On a fresh database that is the normal case, not an error.
	if err := db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(version), -1) FROM schema_migrations`,
	).Scan(&applied); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}

	for i := applied + 1; i < len(migrations); i++ {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", i, err)
		}
		if _, err := tx.ExecContext(ctx, migrations[i]); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %d: %w", i, err)
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (version) VALUES (?)`, i,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %d: %w", i, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", i, err)
		}
	}
	return nil
}
