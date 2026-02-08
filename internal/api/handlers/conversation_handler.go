package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetConversations(c *echo.Context) error {
	return (*c).JSON(http.StatusOK, "Conversations fetched")
}

func GetConversation(c *echo.Context) error {
	return (*c).JSON(http.StatusOK, "Conversation fetched")
}