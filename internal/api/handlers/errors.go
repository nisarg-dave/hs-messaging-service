package handlers

import (
	"errors"
	"net/http"

	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

// writeServiceError maps an error returned from the service layer to a JSON
// response. Validation errors (those wrapping service.ErrValidation) become
// HTTP 400; everything else falls back to 500.
func writeServiceError(c *echo.Context, err error) error {
	if errors.Is(err, service.ErrValidation) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
