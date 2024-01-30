package main

import "github.com/huntermotko/site/internal/server"

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
