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
	// JSONHandler emits one JSON object per line. On ECS, stdout is captured by
	// CloudWatch Logs, which can index structured fields (level, msg, messageId,
	// etc.) for filtering and alerting. Use TextHandler only if you prefer
	// human-readable output during local development.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

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
