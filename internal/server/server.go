package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

// defaultPort is used when the environment supplies no usable port. Without
// it an unset or non-numeric DP/PP silently yields 0, which binds a random
// port — the failure is invisible until something tries to reach the site.
const defaultPort = 8080

// resolvePort reads the port for the current APP_ENV, falling back to
// defaultPort when the variable is missing, empty, or not a number.
func resolvePort() int {
	key := "PP"
	if os.Getenv("APP_ENV") == "development" {
		key = "DP"
	}
	port, err := strconv.Atoi(os.Getenv(key))
	if err != nil || port <= 0 {
		return defaultPort
	}
	return port
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: resolvePort(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
