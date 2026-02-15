package routes

import (
	"hs-messaging-service/internal/api/handlers"

	"github.com/labstack/echo/v5"
)

func RegisterMessageRoutes(e *echo.Echo, messageHandler *handlers.MessageHandler) {
	messageGroup := e.Group("/messages")
	messageGroup.POST("", messageHandler.CreateMessage)
	messageGroup.PATCH("/:id/read", messageHandler.MarkMessageAsRead)
}	