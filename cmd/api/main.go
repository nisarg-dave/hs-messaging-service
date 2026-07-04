package main

import (
	"log/slog"
	"os"

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

	// Composition Root: create shared dependencies once here and pass them down.
	// Pattern: Dependency Injection — same approach used for repositories and
	// services; the logger is not a global singleton.
	//
	// TextHandler produces human-readable lines for local development. Use
	// slog.NewJSONHandler instead in production so log aggregators (Datadog,
	// CloudWatch, Loki, etc.) can parse structured key=value fields automatically.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	db, err := postgres.NewConnection(config)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	messageRepository := postgres.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepository, logger)
	messageHandler := handlers.NewMessageHandler(messageService)

	conversationRepository := postgres.NewConversationRepository(db)
	conversationService := service.NewConversationService(conversationRepository, logger)
	conversationHandler := handlers.NewConversationHandler(conversationService)

	e := echo.New()

	routes.RegisterMessageRoutes(e, messageHandler)
	routes.RegisterConversationRoutes(e, conversationHandler)

	logger.Info("server started", "port", config.ServerPort)
	if err := e.Start(":" + config.ServerPort); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
