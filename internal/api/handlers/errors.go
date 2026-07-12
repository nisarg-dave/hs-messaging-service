package handlers

import (
	"errors"
	"net/http"

	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

// writeServiceError maps an error returned from the service layer to a JSON
// response.
//
//	service.ErrValidation -> 400 Bad Request
//	service.ErrNotFound   -> 404 Not Found
//	anything else         -> 500 Internal Server Error
//
// A "sentinel" here is a fixed, exported error value used for identity checks —
// not a free-form error string. ErrValidation and ErrNotFound are created once
// with errors.New(...) and reused; specific failures wrap them with %w. Like a
// landmark: no matter how many layers wrap the error ("create message:
// validation error: ..."), you can still detect the category by walking the
// chain to that known value. Prefer errors.Is over string matching
// (strings.Contains(err.Error(), "validation")), which breaks if wording changes.
//
// errors.Is(err, target) walks the wrap chain from fmt.Errorf("%w", ...) and
// returns true if target appears anywhere in it. So after the service returns
// "create message: validation error: ...", errors.Is(err, ErrValidation) is
// still true — plain == would fail after wrapping. See also the onion-chain
// docs on ErrValidation in internal/service/errors.go.
//
// err.Error() is the full human-readable string of that chain; we put it in
// the JSON body under "error" while the HTTP status comes from which sentinel
// errors.Is matched.
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
