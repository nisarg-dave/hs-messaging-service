package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// UserIDContextKey is the echo context key under which RequireUserID stores the
// authenticated caller's user ID for downstream handlers.
const UserIDContextKey = "authenticatedUserID"

// RequireUserID returns middleware that authenticates the caller via the X-User-Id header.
//
// Pattern (Middleware / Decorator + Chain of Responsibility): this check runs before
// any handler in the group; unauthorized requests short-circuit without reaching business logic.
//
// SOLID — Single Responsibility: "who is the caller" is owned by this component; handlers
// only map HTTP requests to service calls.
//
// Strategy seam: swapping header auth for JWT later means replacing this middleware only;
// handlers read UserIDFromContext and stay unchanged.
func RequireUserID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := c.Request().Header.Get("X-User-Id")
			if userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "X-User-Id header is required"})
			}
			c.Set(UserIDContextKey, userID)
			return next(c)
		}
	}
}

// UserIDFromContext returns the user ID stored by RequireUserID, or "" if unset or not a string.
func UserIDFromContext(c *echo.Context) string {
	val := c.Get(UserIDContextKey)
	if val == nil {
		return ""
	}
	userID, ok := val.(string)
	if !ok {
		return ""
	}
	return userID
}
