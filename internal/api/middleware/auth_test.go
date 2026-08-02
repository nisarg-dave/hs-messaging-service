package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestRequireUserID_MissingHeader(t *testing.T) {
	e := echo.New()
	nextCalled := false
	handler := RequireUserID()(func(c *echo.Context) error {
		nextCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if nextCalled {
		t.Error("next handler should not be called when header is missing")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	want := map[string]string{"error": "X-User-Id header is required"}
	if body["error"] != want["error"] {
		t.Errorf("body error = %q, want %q", body["error"], want["error"])
	}
}

func TestRequireUserID_HeaderPresent(t *testing.T) {
	e := echo.New()
	var gotUserID string
	handler := RequireUserID()(func(c *echo.Context) error {
		gotUserID = UserIDFromContext(c)
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-Id", "user-42")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if gotUserID != "user-42" {
		t.Errorf("UserIDFromContext = %q, want %q", gotUserID, "user-42")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestUserIDFromContext_NotSet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if got := UserIDFromContext(c); got != "" {
		t.Errorf("UserIDFromContext = %q, want empty string", got)
	}
}
