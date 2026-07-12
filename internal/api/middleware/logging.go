package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

// RequestLogger returns middleware that logs one structured info line per HTTP request.
//
// Pattern: Middleware / Decorator — wraps every handler with cross-cutting request logging
// without modifying the handlers themselves (Chain of Responsibility).
//
// SOLID: Single Responsibility — handlers keep business concerns; request logging lives in one place.
// Open/Closed — new cross-cutting behavior is added by composing middleware, not editing handlers.
//
// Dependency Injection: the logger is passed in from the composition root (main), not a global.
func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)
			_, status := echo.ResolveResponseStatus(c.Response(), err)

			attrs := []any{
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"durationMs", time.Since(start).Milliseconds(),
			}
			if err != nil {
				attrs = append(attrs, "error", err)
			}
			logger.Info("request completed", attrs...)

			return err
		}
	}
}
