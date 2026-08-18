package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/huntermotko/site/internal/database"
)

// newTestServer builds a Server through the same path production uses, so the
// tests exercise real route registration, real template parsing, and the real
// embedded assets. A test that stubs the renderer cannot catch a template that
// fails to parse — which is the failure mode this file most needs to catch,
// since template.Must panics at startup rather than at compile time.
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	s := &Server{cfg: Config{SiteURL: "https://huntermotko.dev"}}
	return s.RegisterRoutes()
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPagesRender(t *testing.T) {
	h := newTestServer(t)

	tests := []struct {
		path     string
		contains []string
	}{
		{
			path: "/",
			contains: []string{
				"Hunter Motko",
				"See the work",
				`<link rel="canonical" href="https://huntermotko.dev/">`,
				`application/ld+json`,
			},
		},
		{
			path: "/work",
			contains: []string{
				"Prebuilt Sheds LLC",
				"Illinois Public Defender Statistics",
				"ilpublicdefenderstats.org",
				// The bench lane must survive on the page, not just in the data.
				"csvq",
				"Unix primitives in C",
				// Year and provenance are the whole reason the tier exists — an
				// undated list cannot show a trajectory, and an unlabelled one
				// lets a reimplementation read as original work. If the template
				// stops rendering them the section quietly loses its point, so
				// assert on the rendered meta line rather than on the data.
				"C · 2026 · original",
				"C · 2024 — 2026 · from spec",
			},
		},
		{
			path: "/about",
			contains: []string{
				"Experience",
				"Education",
				"The Last Mile",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			rec := get(t, h, tc.path)
			if rec.Code != http.StatusOK {
				t.Fatalf("GET %s = %d, want 200", tc.path, rec.Code)
			}
			body := rec.Body.String()
			for _, want := range tc.contains {
				if !strings.Contains(body, want) {
					t.Errorf("GET %s body missing %q", tc.path, want)
				}
			}
		})
	}
}

// TestNoStaleMarketdataCopy guards the specific regression that motivated the
// rewrite: /about described a market-data platform that was abandoned. Shipping
// a confident paragraph about a project that no longer exists is worse than
// shipping nothing, and it is the kind of thing that rots back in silently.
func TestNoStaleMarketdataCopy(t *testing.T) {
	h := newTestServer(t)
	for _, path := range []string{"/", "/about", "/work"} {
		body := strings.ToLower(get(t, h, path).Body.String())
		for _, banned := range []string{"market data platform", "federal economic data", "fred"} {
			if strings.Contains(body, banned) {
				t.Errorf("GET %s still contains abandoned-project copy %q", path, banned)
			}
		}
	}
}

// TestNoExternalScripts guards against reintroducing a CDN dependency. The page
// previously loaded htmx from unpkg.com with no version in the URL, so it
// tracked whatever was published most recently.
func TestNoExternalScripts(t *testing.T) {
	h := newTestServer(t)
	for _, path := range []string{"/", "/about", "/work"} {
		body := get(t, h, path).Body.String()
		for _, host := range []string{"unpkg.com", "cdnjs.cloudflare.com", "cdn.jsdelivr.net", "buymeacoffee.com"} {
			if strings.Contains(body, host) {
				t.Errorf("GET %s references external host %s", path, host)
			}
		}
	}
}

// TestDevOnlyScript keeps the performance observer out of visitors' consoles.
func TestDevOnlyScript(t *testing.T) {
	prod := (&Server{cfg: Config{}}).RegisterRoutes()
	if strings.Contains(get(t, prod, "/").Body.String(), "/js/script.js") {
		t.Error("script.js is served in production, want dev only")
	}

	dev := (&Server{cfg: Config{Dev: true}}).RegisterRoutes()
	if !strings.Contains(get(t, dev, "/").Body.String(), "/js/script.js") {
		t.Error("script.js is missing in development, want it wired up")
	}
}

// TestResumeLinkTracksTheFile: the download link must never appear without the
// route behind it.
func TestResumeLinkTracksTheFile(t *testing.T) {
	h := newTestServer(t) // no ResumePath configured
	if got := get(t, h, "/resume.pdf").Code; got != http.StatusNotFound {
		t.Errorf("GET /resume.pdf = %d with no file configured, want 404", got)
	}
	if strings.Contains(get(t, h, "/").Body.String(), "/resume.pdf") {
		t.Error("home page advertises /resume.pdf while the route is unregistered")
	}
}

// TestResumeRejectsNonFiles: Docker creates an empty *directory* at a bind-mount
// point whose host source is missing, and a directory satisfies a bare
// os.Stat-succeeds check — which would render the download link and then 404.
func TestResumeRejectsNonFiles(t *testing.T) {
	dir := t.TempDir()

	realPDF := filepath.Join(dir, "real.pdf")
	if err := os.WriteFile(realPDF, []byte("%PDF-1.4 not really"), 0o644); err != nil {
		t.Fatal(err)
	}
	emptyPDF := filepath.Join(dir, "empty.pdf")
	if err := os.WriteFile(emptyPDF, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	dirPDF := filepath.Join(dir, "adirectory.pdf")
	if err := os.Mkdir(dirPDF, 0o755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		served bool
	}{
		{name: "real file is served", path: realPDF, served: true},
		{name: "directory is refused", path: dirPDF},
		{name: "empty file is refused", path: emptyPDF},
		{name: "missing path is refused", path: filepath.Join(dir, "nope.pdf")},
		{name: "unset is refused", path: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := (&Server{cfg: Config{ResumePath: tc.path}}).RegisterRoutes()

			wantCode := http.StatusNotFound
			if tc.served {
				wantCode = http.StatusOK
			}
			if got := get(t, h, "/resume.pdf").Code; got != wantCode {
				t.Errorf("GET /resume.pdf = %d, want %d", got, wantCode)
			}

			// The link and the route must never disagree.
			linked := strings.Contains(get(t, h, "/").Body.String(), "/resume.pdf")
			if linked != tc.served {
				t.Errorf("home page links résumé = %v, want %v", linked, tc.served)
			}
		})
	}
}

func TestCanonicalOmittedWithoutSiteURL(t *testing.T) {
	// A canonical built from the Host header would let a forged Host redirect
	// ranking signal to someone else's domain, so no SITE_URL means no tag.
	h := (&Server{cfg: Config{}}).RegisterRoutes()
	body := get(t, h, "/").Body.String()
	if strings.Contains(body, "rel=\"canonical\"") {
		t.Error("canonical tag emitted without SITE_URL set")
	}
	if strings.Contains(body, "og:url") {
		t.Error("og:url emitted without SITE_URL set")
	}
}

func TestSitemapAndRobots(t *testing.T) {
	h := newTestServer(t)

	sitemap := get(t, h, "/sitemap.xml")
	if sitemap.Code != http.StatusOK {
		t.Fatalf("GET /sitemap.xml = %d, want 200", sitemap.Code)
	}
	// Every page in the sitemap must actually serve, or it is a Search Console
	// error rather than a crawling hint.
	for _, path := range []string{"/", "/work", "/about"} {
		loc := "https://huntermotko.dev" + path
		if !strings.Contains(sitemap.Body.String(), "<loc>"+loc+"</loc>") {
			t.Errorf("sitemap missing %s", loc)
		}
		if code := get(t, h, path).Code; code != http.StatusOK {
			t.Errorf("sitemap advertises %s but it returns %d", path, code)
		}
	}
	if strings.Contains(sitemap.Body.String(), "/stats") {
		t.Error("sitemap lists the private /stats page")
	}

	robots := get(t, h, "/robots.txt")
	if !strings.Contains(robots.Body.String(), "Sitemap: https://huntermotko.dev/sitemap.xml") {
		t.Errorf("robots.txt missing absolute sitemap line, got:\n%s", robots.Body.String())
	}
}

// TestHeadIsAllowed: uptime monitors, link checkers, and several link unfurlers
// send HEAD. Echo's e.GET binds GET alone, so every page answered 405 and those
// tools would have read the site as down.
func TestHeadIsAllowed(t *testing.T) {
	h := newTestServer(t)
	for _, path := range []string{"/", "/work", "/about", "/robots.txt", "/sitemap.xml"} {
		req := httptest.NewRequest(http.MethodHead, path, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("HEAD %s = %d, want 200", path, rec.Code)
		}
	}
}

// TestHeadIsNotCounted keeps machine traffic out of the numbers.
func TestHeadIsNotCounted(t *testing.T) {
	db := &countingDB{}
	s := &Server{cfg: Config{StatsToken: "t"}, db: db}
	h := s.RegisterRoutes()

	req := httptest.NewRequest(http.MethodHead, "/", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)
	// The middleware records on a goroutine, so give it a moment to not happen.
	time.Sleep(100 * time.Millisecond)
	if n := db.count(); n != 0 {
		t.Errorf("HEAD recorded %d pageviews, want 0", n)
	}

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	time.Sleep(100 * time.Millisecond)
	if n := db.count(); n != 1 {
		t.Errorf("GET recorded %d pageviews, want 1", n)
	}
}

// countingDB is a database.Service that only counts writes.
type countingDB struct {
	mu sync.Mutex
	n  int
}

func (c *countingDB) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func (c *countingDB) RecordPageview(context.Context, string, int, time.Duration, string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
	return nil
}

func (c *countingDB) Health(context.Context) error { return nil }
func (c *countingDB) Close() error                 { return nil }
func (c *countingDB) PageviewSummary(context.Context, int) (database.Summary, error) {
	return database.Summary{}, nil
}

// TestStatsAcceptsHeaderToken: the header is the form that keeps the token out
// of nginx's access.log, so it has to actually work.
func TestStatsAcceptsHeaderToken(t *testing.T) {
	s := &Server{cfg: Config{StatsToken: "sekrit"}, db: &countingDB{}}
	h := s.RegisterRoutes()

	tests := []struct {
		name   string
		header string
		query  string
		want   int
	}{
		{name: "valid header", header: "sekrit", want: http.StatusOK},
		{name: "valid query", query: "sekrit", want: http.StatusOK},
		{name: "wrong header", header: "nope", want: http.StatusUnauthorized},
		{name: "neither", want: http.StatusUnauthorized},
		// A present-but-wrong header must not fall through to a valid query
		// param, or the header check would be trivially bypassable.
		{name: "wrong header does not fall back", header: "nope", query: "sekrit", want: http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/stats"
			if tc.query != "" {
				url += "?token=" + tc.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tc.header != "" {
				req.Header.Set("X-Stats-Token", tc.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("GET %s (header=%q) = %d, want %d", url, tc.header, rec.Code, tc.want)
			}
		})
	}
}

func TestStatsRequiresToken(t *testing.T) {
	// With no token configured the route must not exist at all, so a host that
	// forgot to set one cannot serve the page to anybody.
	h := (&Server{cfg: Config{}}).RegisterRoutes()
	if got := get(t, h, "/stats").Code; got != http.StatusNotFound {
		t.Errorf("GET /stats = %d with no STATS_TOKEN, want 404", got)
	}
}

func TestJSONLDIsValid(t *testing.T) {
	h := newTestServer(t)
	body := get(t, h, "/").Body.String()

	const open = `<script type="application/ld+json">`
	i := strings.Index(body, open)
	if i < 0 {
		t.Fatal("no JSON-LD block on the home page")
	}
	rest := body[i+len(open):]
	j := strings.Index(rest, "</script>")
	if j < 0 {
		t.Fatal("unterminated JSON-LD block")
	}

	var doc map[string]any
	if err := json.Unmarshal([]byte(rest[:j]), &doc); err != nil {
		t.Fatalf("JSON-LD does not parse: %v", err)
	}
	if doc["@type"] != "Person" {
		t.Errorf("JSON-LD @type = %v, want Person", doc["@type"])
	}
	if doc["name"] != "Hunter Motko" {
		t.Errorf("JSON-LD name = %v", doc["name"])
	}
}

func TestStaticAssetsAreEmbedded(t *testing.T) {
	// These are served out of embed.FS. If the embed directives break, the
	// binary still compiles and every asset 404s at runtime.
	h := newTestServer(t)
	for _, path := range []string{"/css/index.css", "/js/htmx.min.js", "/images/favicon.svg", "/images/headshot.jpeg"} {
		if got := get(t, h, path).Code; got != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, got)
		}
	}
}

// TestFaviconIcoIsServed covers the well-known path. Clients with no HTML in
// front of them — /stats returns JSON, and unfurlers ask directly — never see
// the <link rel="icon"> and fall back here.
//
// The content type is the load-bearing assertion, not the status. middleware.Secure
// sets X-Content-Type-Options: nosniff, so a browser will not rescue a body typed
// from the ".ico" in the URL instead of from the .svg actually being served.
func TestFaviconIcoIsServed(t *testing.T) {
	h := newTestServer(t)

	rec := get(t, h, "/favicon.ico")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /favicon.ico = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "image/svg+xml") {
		t.Errorf("Content-Type = %q, want image/svg+xml", got)
	}
}
