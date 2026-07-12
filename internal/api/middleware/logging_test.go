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

func TestRequestLogger_LogsRequest(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(buf, nil))

	e := echo.New()
	e.Use(RequestLogger(logger))
	e.GET("/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	logLine := buf.String()
	if logLine == "" {
		t.Fatal("expected log output, got empty buffer")
	}

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
	if status, ok := logAttrs["status"].(float64); !ok || int(status) != http.StatusOK {
		t.Errorf("status = %v, want %d", logAttrs["status"], http.StatusOK)
	}
	if _, ok := logAttrs["durationMs"]; !ok {
		t.Errorf("log missing durationMs field: %s", logLine)
	}
}

func TestRequestLogger_ErrorPath(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewJSONHandler(buf, nil))

	e := echo.New()
	handlerErr := errors.New("handler failed")
	handler := RequestLogger(logger)(func(c *echo.Context) error {
		return handlerErr
	})

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler(c)
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
	if status, ok := logAttrs["status"].(float64); !ok || int(status) != http.StatusInternalServerError {
		t.Errorf("status = %v, want %d", logAttrs["status"], http.StatusInternalServerError)
	}
	if errVal, ok := logAttrs["error"].(string); !ok || !strings.Contains(errVal, "handler failed") {
		t.Errorf("error = %v, want message containing %q", logAttrs["error"], "handler failed")
	}
}
