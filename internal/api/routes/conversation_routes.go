package routes

import (
	"hs-messaging-service/internal/api/handlers"
	"hs-messaging-service/internal/api/middleware"

	"github.com/labstack/echo/v5"
)

func RegisterConversationRoutes(e *echo.Echo, conversationHandler *handlers.ConversationHandler) {
	conversationGroup := e.Group("/conversations")
	conversationGroup.Use(middleware.RequireUserID()) // auth applies to every route in the group
	conversationGroup.GET("", conversationHandler.GetConversations)
	conversationGroup.GET("/:userId", conversationHandler.GetConversation)
}
