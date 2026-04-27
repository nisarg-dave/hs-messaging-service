package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hs-messaging-service/internal/domain"

	"github.com/labstack/echo/v5"
)

type fakeConversationService struct {
	listCalledWithUserID string
	listReturn           []domain.ConversationSummary
	listErr              error

	getCalledWithUserID  string
	getCalledWithOtherID string
	getReturn            []domain.Message
	getErr               error
}

func (f *fakeConversationService) ListConversations(userID string) ([]domain.ConversationSummary, error) {
	f.listCalledWithUserID = userID
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listReturn, nil
}

func (f *fakeConversationService) GetConversation(userID, otherID string) ([]domain.Message, error) {
	f.getCalledWithUserID = userID
	f.getCalledWithOtherID = otherID
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.getReturn, nil
}

func newConversationContext(target, userIDHeader string) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if userIDHeader != "" {
		req.Header.Set("X-User-Id", userIDHeader)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestConversationHandler_GetConversations_Success(t *testing.T) {
	now := time.Date(2026, 2, 15, 12, 30, 45, 0, time.UTC)
	fake := &fakeConversationService{
		listReturn: []domain.ConversationSummary{
			{
				UserID:      "user-123",
				LastMessage: domain.LastMessage{Content: "hi", CreatedAt: now},
				UnreadCount: 2,
			},
		},
	}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations", "user-1")
	if err := h.GetConversations(c); err != nil {
		t.Fatalf("GetConversations returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if fake.listCalledWithUserID != "user-1" {
		t.Errorf("service called with userID = %q, want %q", fake.listCalledWithUserID, "user-1")
	}

	var got struct {
		Conversations []domain.ConversationSummary `json:"conversations"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got.Conversations) != 1 || got.Conversations[0].UserID != "user-123" {
		t.Errorf("unexpected conversations payload: %+v", got.Conversations)
	}
	if got.Conversations[0].UnreadCount != 2 {
		t.Errorf("UnreadCount = %d, want 2", got.Conversations[0].UnreadCount)
	}
}

func TestConversationHandler_GetConversations_MissingHeader(t *testing.T) {
	fake := &fakeConversationService{}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations", "")
	if err := h.GetConversations(c); err != nil {
		t.Fatalf("GetConversations returned error: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if fake.listCalledWithUserID != "" {
		t.Errorf("service should not be called, got userID=%q", fake.listCalledWithUserID)
	}
}

func TestConversationHandler_GetConversations_ServiceError(t *testing.T) {
	fake := &fakeConversationService{listErr: errors.New("db down")}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations", "user-1")
	if err := h.GetConversations(c); err != nil {
		t.Fatalf("GetConversations returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestConversationHandler_GetConversation_Success(t *testing.T) {
	fake := &fakeConversationService{
		getReturn: []domain.Message{
			{ID: "m1", SenderID: "user-2", RecipientID: "user-1", Content: "hi", IsRead: true},
		},
	}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations/user-2", "user-1")
	c.SetPathValues(echo.PathValues{{Name: "userId", Value: "user-2"}})

	if err := h.GetConversation(c); err != nil {
		t.Fatalf("GetConversation returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if fake.getCalledWithUserID != "user-1" || fake.getCalledWithOtherID != "user-2" {
		t.Errorf("service called with (%q,%q), want (user-1,user-2)",
			fake.getCalledWithUserID, fake.getCalledWithOtherID)
	}

	var got struct {
		UserID   string           `json:"userId"`
		Messages []domain.Message `json:"messages"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.UserID != "user-2" {
		t.Errorf("userId = %q, want %q", got.UserID, "user-2")
	}
	if len(got.Messages) != 1 || got.Messages[0].ID != "m1" {
		t.Errorf("unexpected messages payload: %+v", got.Messages)
	}
}

func TestConversationHandler_GetConversation_MissingHeader(t *testing.T) {
	fake := &fakeConversationService{}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations/user-2", "")
	c.SetPathValues(echo.PathValues{{Name: "userId", Value: "user-2"}})

	if err := h.GetConversation(c); err != nil {
		t.Fatalf("GetConversation returned error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestConversationHandler_GetConversation_ServiceError(t *testing.T) {
	fake := &fakeConversationService{getErr: errors.New("boom")}
	h := NewConversationHandler(fake)

	c, rec := newConversationContext("/conversations/user-2", "user-1")
	c.SetPathValues(echo.PathValues{{Name: "userId", Value: "user-2"}})

	if err := h.GetConversation(c); err != nil {
		t.Fatalf("GetConversation returned error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
