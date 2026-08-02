package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

// TestRequestLogger_LogsRequest drives a full Echo serve path: register the
// decorator with e.Use, hit a trivial route, then assert the slog JSON line
// written to an in-memory buffer contains the expected fields.
func TestRequestLogger_LogsRequest(t *testing.T) {
	// Buffer-backed JSON handler: *slog.Logger satisfies logging.Logger, and
	// writing to buf lets us inspect the log without printing to stdout.
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(buf, nil))

	e := echo.New()
	// Decorator registration — same call site as cmd/api/main.go.
	e.Use(RequestLogger(logger))
	e.GET("/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	// ServeHTTP runs the full middleware → handler chain.
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	logLine := buf.String()
	if logLine == "" {
		t.Fatal("expected log output, got empty buffer")
	}

	// One JSON object per line from JSONHandler; unmarshal the buffer as a map
	// so we can check structured fields without string matching the whole line.
	var logAttrs map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logAttrs); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}

	if logAttrs["method"] != "GET" {
		t.Errorf("method = %v, want GET", logAttrs["method"])
	}
	if logAttrs["path"] != "/hello" {
		t.Errorf("path = %v, want /hello", logAttrs["path"])
	}
	// JSON numbers decode as float64.
	if status, ok := logAttrs["status"].(float64); !ok || int(status) != http.StatusOK {
		t.Errorf("status = %v, want %d", logAttrs["status"], http.StatusOK)
	}
	if _, ok := logAttrs["durationMs"]; !ok {
		t.Errorf("log missing durationMs field: %s", logLine)
	}
}

// TestRequestLogger_ErrorPath builds the decorator around a handler that
// returns an error (without ServeHTTP) so we can assert both: the error is
// returned unchanged, and the log line still records method/path/status/error.
func TestRequestLogger_ErrorPath(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(buf, nil))

	e := echo.New()
	handlerErr := errors.New("handler failed")
	// RequestLogger(logger) returns MiddlewareFunc; calling it with next yields
	// the decorated HandlerFunc — the Decorator "wrap" step, tested in isolation.
	handler := RequestLogger(logger)(func(c *echo.Context) error {
		return handlerErr
	})

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler(c)
	// Decorator must not swallow errors — Echo's error handling still needs them.
	if !errors.Is(err, handlerErr) {
		t.Fatalf("handler error = %v, want %v", err, handlerErr)
	}

	logLine := buf.String()
	if logLine == "" {
		t.Fatal("expected log output, got empty buffer")
	}

	var logAttrs map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logAttrs); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}

	if logAttrs["method"] != "POST" {
		t.Errorf("method = %v, want POST", logAttrs["method"])
	}
	if logAttrs["path"] != "/fail" {
		t.Errorf("path = %v, want /fail", logAttrs["path"])
	}
	// No status was written before the error; ResolveResponseStatus maps that to 500.
	if status, ok := logAttrs["status"].(float64); !ok || int(status) != http.StatusInternalServerError {
		t.Errorf("status = %v, want %d", logAttrs["status"], http.StatusInternalServerError)
	}
	if errVal, ok := logAttrs["error"].(string); !ok || !strings.Contains(errVal, "handler failed") {
		t.Errorf("error = %v, want message containing %q", logAttrs["error"], "handler failed")
	}
}
