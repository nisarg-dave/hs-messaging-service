package handlers

import (
	"errors"
	"net/http"

	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

// writeServiceError maps an error returned from the service layer to a JSON
// response. It walks the error chain with errors.Is so the wrapped operation
// prefix (e.g. "create message: ...") doesn't hide the underlying sentinel.
//
//	service.ErrValidation -> 400 Bad Request
//	service.ErrNotFound   -> 404 Not Found
//	anything else         -> 500 Internal Server Error
func writeServiceError(c *echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrValidation):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, service.ErrNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
