package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/huntermotko/site/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
  env := os.Getenv("APP_ENV")
  var port int
  if env == "development" {
    port, _ = strconv.Atoi(os.Getenv("DP"))
  } else {
    port, _ = strconv.Atoi(os.Getenv("PP"))
  }

	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
