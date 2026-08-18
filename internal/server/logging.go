package server

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// requestLogger emits at most one line per request, and for most requests none
// at all.
//
// It replaces middleware.Logger(), which logged every request that reached the
// server — every CSS file, every image, every favicon probe. One person reading
// one page produced a line for the page and then a line for each asset on it,
// so the only event actually worth seeing was buried under its own page weight,
// and the journal filled with output nobody reads. A log that is never read
// is not a log.
//
// Two things earn a line:
//
//   - A page view: one line per visit, gated on the same isPageRoute check that
//     recordPageview uses. The log and /stats therefore cannot disagree about
//     what counts as traffic, which matters because the log is the whole record
//     when DB_PATH is unset and metrics are off.
//   - Anything that did not succeed: 4xx and 5xx are logged wherever they come
//     from, static assets included. A missing asset returns 404 and changes
//     nothing a visitor would report, so it is invisible unless it is logged.
//
// 304 Not Modified is excluded despite being over 300. It is the success case
// for a cached asset — a returning visitor generates one per file — so logging
// it would restore precisely the per-asset noise this exists to remove.
//
// middleware.Logger is also deprecated in favour of RequestLogger, so this
// clears that at the same time.
func requestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogStatus:    true,
		LogLatency:   true,
		LogReferer:   true,
		LogError:     true,

		// Run the global error handler before logging, so Status is the code
		// the client actually received. Without this a handler that returns an
		// error is logged as 200, because nothing has written a status yet.
		HandleError: true,

		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			failed := v.Status >= http.StatusMultipleChoices &&
				v.Status != http.StatusNotModified

			// GET only for the success case, matching recordPageview exactly.
			// The page routes also answer HEAD for uptime monitors and link
			// checkers, and a monitor polling once a minute is ~1,400 lines a
			// day — the same burial this change exists to undo, just from a
			// different direction. A failing HEAD still logs, via failed.
			isView := isPageRoute(v.RoutePath) && v.Method == http.MethodGet

			if !failed && !isView {
				return nil
			}

			// URIPath, not URI: the query string is dropped on purpose. /stats
			// accepts its token as a query parameter, and a log line is a far
			// more durable place for a secret to sit than the browser history
			// that form already trades away.
			lat := v.Latency.Round(time.Microsecond)

			switch {
			case v.Error != nil:
				log.Printf("%s %s %d %s err=%v",
					v.Method, v.URIPath, v.Status, lat, v.Error)
			case v.Referer != "":
				// %q on the referrer, which is attacker-controlled: a raw
				// value containing a newline would otherwise forge additional
				// log lines below the real one.
				log.Printf("%s %s %d %s ref=%q",
					v.Method, v.URIPath, v.Status, lat, truncate(v.Referer, 300))
			default:
				log.Printf("%s %s %d %s", v.Method, v.URIPath, v.Status, lat)
			}
			return nil
		},
	})
}
