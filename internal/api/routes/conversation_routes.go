package routes

import (
	"hs-messaging-service/internal/api/handlers"

	"github.com/labstack/echo/v5"
)

func RegisterConversationRoutes(e *echo.Echo, conversationHandler *handlers.ConversationHandler) {
	conversationGroup := e.Group("/conversations")
	conversationGroup.GET("", conversationHandler.GetConversations)
	conversationGroup.GET("/:userId", conversationHandler.GetConversation)
}
