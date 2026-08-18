package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// captureLog redirects the standard logger for the duration of one test.
//
// No test in this package calls t.Parallel, so swapping a process-global is
// safe here; the restore is deferred by the caller through t.Cleanup.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	flags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(flags)
	})
	return &buf
}

// TestLogIsOnePerPageView pins the whole point of the request logger: a visit
// produces one line, and the assets that visit drags in produce none. Asserting
// the count rather than the content is deliberate — the failure being guarded
// against is not a wrong line, it is dozens of right ones.
func TestLogIsOnePerPageView(t *testing.T) {
	h := newTestServer(t)
	buf := captureLog(t)

	for _, path := range []string{
		"/",
		"/css/index.css",
		"/js/htmx.min.js",
		"/images/favicon.svg",
		"/robots.txt",
		"/sitemap.xml",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		h.ServeHTTP(httptest.NewRecorder(), req)
	}

	lines := logLines(buf)
	if len(lines) != 1 {
		t.Fatalf("want 1 line for 1 page view, got %d:\n%s",
			len(lines), strings.Join(lines, "\n"))
	}
	if !strings.HasPrefix(lines[0], "GET / 200 ") {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

// TestLogSkipsHead mirrors TestHeadIsNotCounted. An uptime monitor probing once
// a minute is ~1,400 requests a day, and logging them buries the handful of
// real visits the log exists to show.
func TestLogSkipsHead(t *testing.T) {
	h := newTestServer(t)
	buf := captureLog(t)

	req := httptest.NewRequest(http.MethodHead, "/", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if lines := logLines(buf); len(lines) != 0 {
		t.Errorf("HEAD should not log, got:\n%s", strings.Join(lines, "\n"))
	}
}

// TestLogRecordsFailures covers the other half: a request that fails is logged
// wherever it came from, including the static tree, because a 404 on an asset
// changes nothing a visitor would think to report.
func TestLogRecordsFailures(t *testing.T) {
	h := newTestServer(t)
	buf := captureLog(t)

	for _, path := range []string{"/css/missing.css", "/wp-login.php"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		h.ServeHTTP(httptest.NewRecorder(), req)
	}

	lines := logLines(buf)
	if len(lines) != 2 {
		t.Fatalf("want 2 failure lines, got %d:\n%s",
			len(lines), strings.Join(lines, "\n"))
	}
	for _, line := range lines {
		if !strings.Contains(line, " 404 ") {
			t.Errorf("expected a 404 in %q", line)
		}
	}
}

// TestLogOmitsQueryString guards the /stats token. It is accepted as a query
// parameter, and a log line outlives the browser history that form already
// trades away — so the path is logged and the query is not.
//
// The request is driven to 401 on purpose, because that is the case that
// actually reaches the logger: a successful /stats call is not a page route and
// produces no line at all, so it could not catch a regression here.
func TestLogOmitsQueryString(t *testing.T) {
	s := &Server{cfg: Config{StatsToken: "the-real-token"}, db: &countingDB{}}
	h := s.RegisterRoutes()
	buf := captureLog(t)

	req := httptest.NewRequest(http.MethodGet, "/stats?token=guessed-secret", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("stats: want 401, got %d", rec.Code)
	}
	got := buf.String()
	if !strings.Contains(got, "/stats") {
		t.Fatalf("expected the failed request to be logged, got:\n%s", got)
	}
	if strings.Contains(got, "token=") || strings.Contains(got, "guessed-secret") {
		t.Errorf("query string leaked into the log:\n%s", got)
	}
}

func logLines(buf *bytes.Buffer) []string {
	out := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(out) == 1 && out[0] == "" {
		return nil
	}
	return out
}
