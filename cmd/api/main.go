package main

import (
	"log"

	"hs-messaging-service/internal/config"
	"hs-messaging-service/internal/repository/postgres"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config := config.Load()
	_, err := postgres.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}