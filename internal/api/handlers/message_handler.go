package handlers

import (
	"net/http"

	"hs-messaging-service/internal/domain"

	"github.com/labstack/echo/v5"
)

// MessageService is the behavior MessageHandler needs.
// Defining the interface here (the consumer) lets us swap in fakes for tests
// while production code keeps passing in *service.MessageService.
type MessageService interface {
	CreateMessage(message *domain.Message) error
	MarkMessageAsRead(messageID string) (*domain.Message, error)
}

type MessageHandler struct {
	messageService MessageService
}

func NewMessageHandler(messageService MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(c *echo.Context) error {
	message := new(domain.Message)

	if err := c.Bind(message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.messageService.CreateMessage(message); err != nil {
		return writeServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) MarkMessageAsRead(c *echo.Context) error {
	messageID := c.Param("id")
	message, err := h.messageService.MarkMessageAsRead(messageID)
	if err != nil {
		return writeServiceError(c, err)
	}
	return c.JSON(http.StatusOK, message)
}
