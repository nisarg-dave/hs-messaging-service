package main

import (
	"hs-messaging-service/internal/config"
	"hs-messaging-service/internal/repository/postgres"
	"log"
)

func main() {
	config := config.Load()
	_, err := postgres.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}