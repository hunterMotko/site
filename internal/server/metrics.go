package server

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// secureEqual compares two secrets in constant time. A plain == returns as soon
// as it finds a differing byte, which leaks the length of the matching prefix
// through timing and makes the token recoverable one byte at a time.
func secureEqual(got, want string) bool {
	if want == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

// recordPageview writes one row per page request.
//
// Only page routes are recorded — static assets are skipped, because a "view"
// counted once per CSS file and once per image tells you about the page's asset
// count rather than about its traffic.
//
// Recording happens after the handler has run, on a background goroutine with
// its own context. Doing it inline would put a disk write on the response path
// of every request, and using the request's context would cancel the write the
// moment the client disconnected — which is exactly when a slow or abandoned
// request is most worth knowing about.
func (s *Server) recordPageview(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)

		path := c.Path()
		if !isPageRoute(path) {
			return err
		}

		// GET only. The page routes also answer HEAD, for uptime monitors and
		// link checkers — and an uptime monitor polling every 60 seconds is
		// ~1,400 requests a day, which would bury the handful of real visits
		// this is meant to measure under machine traffic.
		if c.Request().Method != http.MethodGet {
			return err
		}

		status := c.Response().Status
		// Echo sets Response().Status only once something is written; an error
		// handled upstream may leave it at the zero value.
		if status == 0 {
			status = http.StatusOK
		}

		go func(path string, status int, dur time.Duration, referrer string) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := s.db.RecordPageview(ctx, path, status, dur, referrer); err != nil {
				log.Printf("record pageview %s: %v", path, err)
			}
		}(path, status, time.Since(start), truncate(c.Request().Referer(), 300))

		return err
	}
}

// isPageRoute reports whether a matched route is a page worth counting. It
// tests the route pattern from the router, not the raw URL, so a request for a
// nonexistent asset cannot inflate the numbers.
func isPageRoute(path string) bool {
	switch path {
	case "/", "/work", "/about":
		return true
	}
	return false
}

// truncate bounds a stored string. Referer is attacker-controlled and
// unbounded; without a limit a single request could write megabytes into the
// database.
func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) > n {
		return s[:n]
	}
	return s
}

// Stats renders the request counts. It is registered only when STATS_TOKEN is
// set, and requires that token.
//
// Kept private on purpose: a public counter on a site whose whole argument is
// "I am worth hiring" reads badly at the traffic levels a personal site
// actually gets, and the numbers are for deciding whether outreach is working,
// not for visitors.
func (s *Server) Stats(c echo.Context) error {
	// The header is the preferred form. A query string lands in nginx's
	// access.log in plaintext and stays there through log rotation, and in
	// browser history besides. The query parameter is kept because it is the
	// only way to open this in a browser, and the blast radius is small — the
	// token guards a pageview counter, not credentials.
	//
	// Constant-time comparison so the token cannot be recovered a byte at a
	// time by measuring how long the failure takes.
	token := c.Request().Header.Get("X-Stats-Token")
	if token == "" {
		token = c.QueryParam("token")
	}
	if !secureEqual(token, s.cfg.StatsToken) {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	summary, err := s.db.PageviewSummary(ctx, 30)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "stats unavailable")
	}
	return c.JSON(http.StatusOK, summary)
}
