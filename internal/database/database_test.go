package database

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestDB(t *testing.T) Service {
	t.Helper()
	// A file in the test's temp dir rather than :memory:. With
	// SetMaxOpenConns(1) an in-memory database would work, but a file exercises
	// the same open path production uses, including migrations against a real
	// empty file.
	db, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestNewRunsMigrationsAndIsHealthy(t *testing.T) {
	db := newTestDB(t)
	if err := db.Health(context.Background()); err != nil {
		t.Fatalf("Health on a fresh database: %v", err)
	}
}

// TestMigrationsAreIdempotent covers the reopen path: every deploy restarts the
// process against a database that already has the schema.
func TestMigrationsAreIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	first, err := New(path)
	if err != nil {
		t.Fatalf("first New: %v", err)
	}
	ctx := context.Background()
	if err := first.RecordPageview(ctx, "/work", 200, 12*time.Millisecond, ""); err != nil {
		t.Fatalf("RecordPageview: %v", err)
	}
	first.Close()

	second, err := New(path)
	if err != nil {
		t.Fatalf("reopening an already-migrated database: %v", err)
	}
	defer second.Close()

	// The row written before the restart must survive it — a migration that
	// re-ran destructively would be invisible except for the missing data.
	summary, err := second.PageviewSummary(ctx, 30)
	if err != nil {
		t.Fatalf("PageviewSummary: %v", err)
	}
	if summary.TotalViews != 1 {
		t.Errorf("TotalViews after reopen = %d, want 1", summary.TotalViews)
	}
}

func TestPageviewSummary(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := db.RecordPageview(ctx, "/work", 200, 10*time.Millisecond, "https://news.ycombinator.com/"); err != nil {
			t.Fatalf("RecordPageview: %v", err)
		}
	}
	if err := db.RecordPageview(ctx, "/about", 200, 5*time.Millisecond, ""); err != nil {
		t.Fatalf("RecordPageview: %v", err)
	}

	summary, err := db.PageviewSummary(ctx, 30)
	if err != nil {
		t.Fatalf("PageviewSummary: %v", err)
	}

	if summary.TotalViews != 4 {
		t.Errorf("TotalViews = %d, want 4", summary.TotalViews)
	}
	if len(summary.ByPath) == 0 {
		t.Fatal("ByPath is empty")
	}
	// Ordered by volume, so the busiest page leads.
	if summary.ByPath[0].Path != "/work" || summary.ByPath[0].Views != 3 {
		t.Errorf("ByPath[0] = %+v, want /work with 3 views", summary.ByPath[0])
	}

	// Blank referrers are excluded: a row per direct visit would drown the
	// field that answers whether a shared link brought anyone.
	if len(summary.ByReferrer) != 1 {
		t.Fatalf("ByReferrer has %d entries, want 1", len(summary.ByReferrer))
	}
	if summary.ByReferrer[0].Views != 3 {
		t.Errorf("referrer views = %d, want 3", summary.ByReferrer[0].Views)
	}
}

// TestSummaryCountsEveryRequest guards the percentile query's row count.
//
// The bucket filter was originally a WHERE, which dropped the slowest 5% of
// requests from COUNT(*) as well as from the percentile — so every path
// reported 95% of its real traffic. It needs more than 20 rows on one path to
// reproduce, which is why the earlier four-row test passed straight through it.
func TestSummaryCountsEveryRequest(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	const n = 100
	for i := 0; i < n; i++ {
		// Ascending durations, so the slow tail lands in the top bucket.
		if err := db.RecordPageview(ctx, "/work", 200, time.Duration(i)*time.Microsecond, ""); err != nil {
			t.Fatalf("RecordPageview: %v", err)
		}
	}

	summary, err := db.PageviewSummary(ctx, 30)
	if err != nil {
		t.Fatalf("PageviewSummary: %v", err)
	}
	if summary.TotalViews != n {
		t.Errorf("TotalViews = %d, want %d", summary.TotalViews, n)
	}
	if len(summary.ByPath) != 1 {
		t.Fatalf("ByPath has %d entries, want 1", len(summary.ByPath))
	}
	if got := summary.ByPath[0].Views; got != n {
		t.Errorf("ByPath views = %d, want %d — the slow tail is being dropped from the count", got, n)
	}
	// Durations run 0..99µs, so the 95th percentile sits near 95µs = 0.095ms.
	// The assertion that matters is that it is not zero, which is what whole
	// milliseconds produced for every sub-millisecond render.
	if summary.ByPath[0].P95MS <= 0 {
		t.Errorf("P95MS = %v, want a non-zero sub-millisecond value", summary.ByPath[0].P95MS)
	}
}

// TestEmptySummaryMarshalsAsArrays: no traffic must not change the JSON shape.
func TestEmptySummaryMarshalsAsArrays(t *testing.T) {
	db := newTestDB(t)
	summary, err := db.PageviewSummary(context.Background(), 30)
	if err != nil {
		t.Fatalf("PageviewSummary: %v", err)
	}
	b, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, want := range []string{`"by_path":[]`, `"by_referrer":[]`} {
		if !strings.Contains(string(b), want) {
			t.Errorf("empty summary JSON missing %s, got %s", want, b)
		}
	}
}

// TestNoVisitorIdentifierColumns is a schema assertion, not a behaviour test.
// The privacy claim on this site is "these columns do not exist", and that is
// only true for as long as nobody adds one.
func TestNoVisitorIdentifierColumns(t *testing.T) {
	db := newTestDB(t)
	svc, ok := db.(*service)
	if !ok {
		t.Fatal("expected *service")
	}

	rows, err := svc.db.QueryContext(context.Background(), `PRAGMA table_info(pageviews)`)
	if err != nil {
		t.Fatalf("table_info: %v", err)
	}
	defer rows.Close()

	banned := map[string]bool{
		"ip": true, "ip_address": true, "remote_addr": true,
		"user_agent": true, "ua": true, "session_id": true, "visitor_id": true,
	}
	for rows.Next() {
		var (
			cid, notnull, pk int
			name, ctype      string
			dflt             any
		)
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		if banned[name] {
			t.Errorf("pageviews has visitor-identifying column %q", name)
		}
	}
}
