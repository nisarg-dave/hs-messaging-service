package routes

import (
	"hs-messaging-service/internal/api/handlers"

	"github.com/labstack/echo/v5"
)

func RegisterConversationRoutes(e *echo.Echo) {
	conversationGroup := e.Group("/conversations")
	conversationGroup.GET("", handlers.GetConversations)
	conversationGroup.GET("/:id", handlers.GetConversation)
}