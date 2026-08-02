package middleware

import (
	"time"

	"hs-messaging-service/internal/logging"

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
// Dependency Injection: logging.Logger is passed in from the composition root (main), not a global.
// The interface lives in internal/logging so both middleware and services share one contract
// without either layer importing the other.
func RequestLogger(logger logging.Logger) echo.MiddlewareFunc {
	// Echo middleware is a function that takes the next HandlerFunc and returns
	// a new HandlerFunc — same type, so this is a Decorator: wrap, add behavior,
	// then call next. Outer return type matches what e.Use(...) expects.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Capture start before calling the rest of the chain so duration
			// includes the handler (and any later middleware).
			start := time.Now()

			// Delegate to the next HandlerFunc (another middleware or the route
			// handler). We do not swallow the error — return it below so Echo
			// can still apply its error handling.
			err := next(c)

			// Echo v5's Response() is a plain http.ResponseWriter; ResolveResponseStatus
			// reads the status Echo wrote, or maps a returned error to 500 when
			// no status was set yet.
			_, status := echo.ResolveResponseStatus(c.Response(), err)

			attrs := []any{
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"durationMs", time.Since(start).Milliseconds(),
			}
			if err != nil {
				// Include the error string as a structured field for CloudWatch filters.
				attrs = append(attrs, "error", err)
			}
			logger.Info("request completed", attrs...)

			// Propagate the original error unchanged so outer middleware / Echo
			// still see it.
			return err
		}
	}
}
