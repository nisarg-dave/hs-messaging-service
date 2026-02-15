package main

import (
	"log"

	"hs-messaging-service/internal/api/handlers"
	"hs-messaging-service/internal/api/routes"
	"hs-messaging-service/internal/config"
	"hs-messaging-service/internal/repository/postgres"
	"hs-messaging-service/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	godotenv.Load()
	config := config.Load()
	db, err := postgres.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	messageRepository := postgres.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepository)
	messageHandler := handlers.NewMessageHandler(messageService)

	e := echo.New()

	routes.RegisterMessageRoutes(e, messageHandler)
	
	log.Printf("Server started on port %s", config.ServerPort)
	if err := e.Start(":" + config.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}