package handlers

import (
	"net/http"

	"hs-messaging-service/internal/domain"

	"github.com/labstack/echo/v5"
)

type ConversationService interface {
	ListConversations(userID string) ([]domain.ConversationSummary, error)
	GetConversation(userID, otherID string) ([]domain.Message, error)
}

type ConversationHandler struct {
	conversationService ConversationService
}

func NewConversationHandler(conversationService ConversationService) *ConversationHandler {
	return &ConversationHandler{conversationService: conversationService}
}

func (h *ConversationHandler) GetConversations(c *echo.Context) error {
	userID := c.Request().Header.Get("X-User-Id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "X-User-Id header is required"})
	}

	conversations, err := h.conversationService.ListConversations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"conversations": conversations})
}

func (h *ConversationHandler) GetConversation(c *echo.Context) error {
	userID := c.Request().Header.Get("X-User-Id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "X-User-Id header is required"})
	}

	otherID := c.Param("userId")
	messages, err := h.conversationService.GetConversation(userID, otherID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"userId":   otherID,
		"messages": messages,
	})
}
