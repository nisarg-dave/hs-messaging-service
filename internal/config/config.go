package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort string
}

func Load() *Config {
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")	
	serverPort := os.Getenv("SERVER_PORT")

	databaseURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, dbUser, dbPassword, dbName, dbPort)
	return &Config{DatabaseURL: databaseURL, ServerPort: serverPort}
}