package database

import (
	"context"
	"fmt"
	"math"
	"time"
)

// PathCount is one row of the per-path breakdown.
type PathCount struct {
	Path  string `json:"path"`
	Views int    `json:"views"`
	// P95MS is the 95th-percentile server render time in milliseconds. The mean
	// hides the requests that were actually slow, which are the only ones worth
	// acting on. Reported as a float because these renders are routinely under
	// a millisecond and an integer would report every one of them as zero.
	P95MS float64 `json:"p95_ms"`
}

// ReferrerCount is one row of the referrer breakdown — the field that answers
// whether a link shared somewhere actually brought anyone.
type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Views    int    `json:"views"`
}

// Summary is the /stats payload.
type Summary struct {
	Days        int             `json:"days"`
	TotalViews  int             `json:"total_views"`
	ByPath      []PathCount     `json:"by_path"`
	ByReferrer  []ReferrerCount `json:"by_referrer"`
	GeneratedAt time.Time       `json:"generated_at"`
}

// RecordPageview stores a single page request.
func (s *service) RecordPageview(
	ctx context.Context,
	path string,
	status int,
	dur time.Duration,
	referrer string,
) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO pageviews (path, status, duration_us, referrer)
		 VALUES (?, ?, ?, ?)`,
		path, status, dur.Microseconds(), referrer,
	)
	if err != nil {
		return fmt.Errorf("insert pageview: %w", err)
	}
	return nil
}

// PageviewSummary aggregates the last n days.
func (s *service) PageviewSummary(ctx context.Context, days int) (Summary, error) {
	if days <= 0 {
		days = 30
	}
	// Bound the window rather than interpolating the caller's number into SQL.
	// days is an int from our own code today, but the parameter is the kind
	// that grows a query-string source later.
	since := fmt.Sprintf("-%d days", days)

	// Slices start empty rather than nil so the JSON carries [] instead of null.
	// A consumer that iterates the result should not have to special-case "no
	// traffic yet" as a different type.
	out := Summary{
		Days:        days,
		GeneratedAt: time.Now().UTC(),
		ByPath:      []PathCount{},
		ByReferrer:  []ReferrerCount{},
	}

	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pageviews WHERE created_at >= datetime('now', ?)`,
		since,
	).Scan(&out.TotalViews); err != nil {
		return Summary{}, fmt.Errorf("total views: %w", err)
	}

	// SQLite has no percentile function, so the p95 is taken positionally:
	// NTILE(20) splits each path's durations into twentieths and the top of the
	// 19th is the 95th percentile.
	//
	// The bucket filter is a CASE inside the aggregate, not a WHERE. As a WHERE
	// it also removed those rows from COUNT(*), so every path reported 95% of
	// its actual traffic — a quiet undercount that looks entirely plausible.
	//
	// The outer COALESCE covers a path whose rows all land in bucket 20, which
	// happens whenever it has fewer than 20 rows in the window; falling back to
	// the overall MAX is right there, since with that few samples the maximum is
	// the best percentile estimate available.
	rows, err := s.db.QueryContext(ctx,
		`SELECT path,
		        COUNT(*) AS views,
		        COALESCE(
		          MAX(CASE WHEN bucket <= 19 THEN duration_us END),
		          MAX(duration_us),
		          0
		        ) AS p95_us
		   FROM (
		     SELECT path, duration_us,
		            NTILE(20) OVER (PARTITION BY path ORDER BY duration_us) AS bucket
		       FROM pageviews
		      WHERE created_at >= datetime('now', ?)
		   )
		  GROUP BY path
		  ORDER BY views DESC`,
		since,
	)
	if err != nil {
		return Summary{}, fmt.Errorf("by path: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			pc    PathCount
			p95us int64
		)
		if err := rows.Scan(&pc.Path, &pc.Views, &p95us); err != nil {
			return Summary{}, fmt.Errorf("scan path row: %w", err)
		}
		// Rounded to two decimals: sub-microsecond precision is noise, and the
		// full float prints as an unreadable 0.4830000000000001.
		pc.P95MS = math.Round(float64(p95us)/10) / 100
		out.ByPath = append(out.ByPath, pc)
	}
	if err := rows.Err(); err != nil {
		return Summary{}, fmt.Errorf("iterate path rows: %w", err)
	}

	refRows, err := s.db.QueryContext(ctx,
		`SELECT referrer, COUNT(*) AS views
		   FROM pageviews
		  WHERE created_at >= datetime('now', ?)
		    AND referrer <> ''
		  GROUP BY referrer
		  ORDER BY views DESC
		  LIMIT 25`,
		since,
	)
	if err != nil {
		return Summary{}, fmt.Errorf("by referrer: %w", err)
	}
	defer refRows.Close()
	for refRows.Next() {
		var rc ReferrerCount
		if err := refRows.Scan(&rc.Referrer, &rc.Views); err != nil {
			return Summary{}, fmt.Errorf("scan referrer row: %w", err)
		}
		out.ByReferrer = append(out.ByReferrer, rc)
	}
	if err := refRows.Err(); err != nil {
		return Summary{}, fmt.Errorf("iterate referrer rows: %w", err)
	}

	return out, nil
}
