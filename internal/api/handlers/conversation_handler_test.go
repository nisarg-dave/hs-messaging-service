package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestConversationHandlers(t *testing.T) {
	tests := []struct {
		name    string
		handler func(*echo.Context) error
		want    string
	}{
		{
			name:    "get conversations",
			handler: GetConversations,
			want:    "Conversations fetched",
		},
		{
			name:    "get conversation",
			handler: GetConversation,
			want:    "Conversation fetched",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := tt.handler(c); err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != http.StatusOK {
				t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
			}

			var got string
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decode response body: %v", err)
			}

			if got != tt.want {
				t.Errorf("response body = %q, want %q", got, tt.want)
			}
		})
	}
}
