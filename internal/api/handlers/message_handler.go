package handlers

import (
	"net/http"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

// MessageService is the behavior MessageHandler needs.
// Defining the interface here (the consumer) lets us swap in fakes for tests
// while production code keeps passing in *service.MessageService.
type MessageService interface {
	CreateMessage(req *service.CreateMessageRequest) (*domain.Message, error)
	MarkMessageAsRead(messageID string) (*domain.Message, error)
}

type MessageHandler struct {
	messageService MessageService
}

func NewMessageHandler(messageService MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(c *echo.Context) error {
	// Bind into the service request DTO, not domain.Message, so the client
	// can't supply ID / IsRead / CreatedAt / UpdatedAt and have them persisted.
	req := new(service.CreateMessageRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message, err := h.messageService.CreateMessage(req)
	if err != nil {
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
