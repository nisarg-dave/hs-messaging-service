package handlers

import (
	"net/http"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(c *echo.Context) error {
	// Creates a pointer to an empty message struct
	message := new(domain.Message)

	// Bind reads the JSON from the http request body and converts it to a message struct
	// Echo context represents the http request and response
	err := c.Bind(message)
	if err != nil {
		// err.Error() returns error message as a string
		// can't convert err directly to a string via JSON
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.messageService.CreateMessage(message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Message created")
}