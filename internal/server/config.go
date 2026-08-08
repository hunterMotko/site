package server

import (
	"os"
	"strconv"
	"strings"
)

// defaultPort is used when the environment supplies no usable port. Without
// it an unset or non-numeric DP/PP silently yields 0, which binds a random
// port — the failure is invisible until something tries to reach the site.
const defaultPort = 8080

// Config is every environment-driven setting, resolved once at startup so the
// rest of the package reads fields instead of calling os.Getenv from wherever
// it happens to need a value. Settings scattered through handlers are settings
// nobody can enumerate.
type Config struct {
	// Port to bind. See resolvePort for how it is chosen.
	Port int

	// Dev is true in development. It gates the performance-observer script,
	// which belongs in your console and not in a visitor's.
	Dev bool

	// SiteURL is the absolute origin, e.g. "https://huntermotko.dev".
	//
	// canonical and og:url are built from this and never from the request's
	// Host header. A forged Host reflected into a canonical tag hands search
	// engines someone else's URL as the authority for your content, and the
	// tag ships inside a page real visitors fetch. robots.txt and sitemap.xml
	// may fall back to Host (see origin) because those are fetched by the
	// crawler itself, which sends the true value.
	SiteURL string

	// DBPath is the SQLite file backing request metrics. Empty disables
	// metrics entirely rather than failing startup — the site's job is
	// serving pages, and it should still do that with no database present.
	DBPath string

	// StatsToken guards /stats. Empty leaves the route unregistered, so the
	// page cannot be reachable by accident on a host that never set it.
	StatsToken string

	// ResumePath is the file served at /resume.pdf. The route registers only
	// when the file is actually present, so the download link is never dead.
	ResumePath string
}

// LoadConfig reads configuration from the environment.
func LoadConfig() Config {
	dev := os.Getenv("APP_ENV") == "development"
	return Config{
		Port:       resolvePort(dev),
		Dev:        dev,
		SiteURL:    strings.TrimRight(os.Getenv("SITE_URL"), "/"),
		DBPath:     os.Getenv("DB_PATH"),
		StatsToken: os.Getenv("STATS_TOKEN"),
		ResumePath: os.Getenv("RESUME_PATH"),
	}
}

// resolvePort reads the port for the current APP_ENV, falling back to
// defaultPort when the variable is missing, empty, or not a number.
//
// Note that in production APP_ENV is not "development", so this reads PP — and
// binding a port below 1024 requires root or a reverse proxy in front. nginx
// currently terminates TLS on the droplet, so PP should name the port nginx
// proxies to, not 80.
func resolvePort(dev bool) int {
	key := "PP"
	if dev {
		key = "DP"
	}
	port, err := strconv.Atoi(os.Getenv(key))
	if err != nil || port <= 0 || port > 65535 {
		return defaultPort
	}
	return port
}
