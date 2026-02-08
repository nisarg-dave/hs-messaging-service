package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func CreateMessage(c *echo.Context) error {
	return (*c).JSON(http.StatusCreated, "Message created")
}

func MarkMessageAsRead(c *echo.Context) error {
	return (*c).JSON(http.StatusOK, "Message marked as read")
}