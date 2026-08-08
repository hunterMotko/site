package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/huntermotko/site/internal/server"
)

func main() {
	srv := server.NewServer()

	log.Printf("listening on %s", srv.Addr)

	// The previous version panicked with the fixed string "cannot start
	// server", which discarded the only useful part: whether the port was
	// taken, the address was malformed, or the process lacked permission to
	// bind. ErrServerClosed is the normal outcome of a graceful shutdown and
	// is not a failure.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server on %s: %v", srv.Addr, err)
	}
}
