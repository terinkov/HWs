package main

import (
	"flag"
	"log"

	"github.com/terinkov_HW2/storage"
	"github.com/terinkov_HW2/server"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {
	addr := flag.String("addr", ":8080", "address for http server")

	taskStorage := storage.NewRamStorage()
	sessionStorage := storage.NewRamSessionRepository()
	userStorage := storage.NewRamUserRepository()

	log.Printf("Starting server on %s", *addr)
	if err := server.CreateAndRunServer(taskStorage,sessionStorage, userStorage, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
