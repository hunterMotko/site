package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/huntermotko/site/internal/database"
	"github.com/joho/godotenv"
)

// Server holds the wired dependencies the handlers need.
type Server struct {
	cfg Config
	db  database.Service

	// hasResume is set during route registration if the résumé file exists,
	// and read by newMeta so the contact block and the route agree.
	hasResume bool
}

func init() {
	// Precedence, highest first: the real environment, then .env.local, then
	// .env. godotenv.Load never overwrites a variable that is already set, and
	// for keys defined in several files the first file listed wins — so this
	// ordering falls out of the argument order.
	//
	// This replaces godotenv/autoload, which read only .env. That file is
	// tracked and has every value empty, while the real values live in
	// .env.local, which is gitignored and was read by nothing at all — a split
	// that stayed invisible while the port was hardcoded.
	//
	// Overload would be wrong here: it overwrites the real environment, so a
	// stale .env.local left in a deploy directory would silently beat the
	// variables systemd or Docker actually set.
	//
	// Missing files are the normal case in production, so the error is ignored.
	_ = godotenv.Load(".env.local", ".env")
}

// NewServer builds the HTTP server. Storage is optional: if no DB_PATH is set,
// or the database cannot be opened, the site still serves every page — it just
// records no metrics. Losing request counts is not a reason to take the site
// down.
func NewServer() *http.Server {
	cfg := LoadConfig()

	// Both metrics and /stats degrade to "off" rather than failing startup, so
	// each disabled path says so and says why. Silence here is what makes a
	// missing DB_PATH indistinguishable from a bug: the route simply is not
	// registered, and the only symptom is a bare 404 from the catch-all.
	var db database.Service
	switch {
	case cfg.DBPath == "":
		log.Print("metrics: disabled (DB_PATH is not set) — pages still serve normally")
	default:
		var err error
		db, err = database.New(cfg.DBPath)
		if err != nil {
			log.Printf("metrics: disabled (opening %s failed: %v)", cfg.DBPath, err)
			db = nil
		} else {
			log.Printf("metrics: recording to %s", cfg.DBPath)
		}
	}

	switch {
	case db == nil:
		log.Print("/stats: not registered (metrics are disabled) — requests will 404")
	case cfg.StatsToken == "":
		log.Print("/stats: not registered (STATS_TOKEN is not set) — requests will 404")
	default:
		log.Print("/stats: registered — requires ?token=")
	}

	s := &Server{cfg: cfg, db: db}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
