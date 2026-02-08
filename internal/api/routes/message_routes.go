package routes

import (
	"hs-messaging-service/internal/api/handlers"

	"github.com/labstack/echo/v5"
)

func RegisterMessageRoutes(e *echo.Echo) {
	messageGroup := e.Group("/messages")
	messageGroup.POST("", handlers.CreateMessage)
	messageGroup.PUT("/:id/read", handlers.MarkMessageAsRead)
}

func RegisterConversationRoutes(e *echo.Echo) {
	conversationGroup := e.Group("/conversations")
	conversationGroup.GET("", handlers.GetConversations)
	conversationGroup.GET("/:id", handlers.GetConversation)
}