package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/service"

	"github.com/labstack/echo/v5"
)

// fakeMessageService satisfies the MessageService interface so the handler
// accepts it. Each field lets a test control what the fake returns.
type fakeMessageService struct {
	createCalledWith *service.CreateMessageRequest
	createReturn     *domain.Message
	createErr        error

	markCalledWithID string
	markReturnMsg    *domain.Message
	markErr          error
}

func (f *fakeMessageService) CreateMessage(req *service.CreateMessageRequest) (*domain.Message, error) {
	f.createCalledWith = req
	if f.createErr != nil {
		return nil, f.createErr
	}
	if f.createReturn != nil {
		return f.createReturn, nil
	}
	return &domain.Message{
		ID:          "fake-id",
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
		Content:     req.Content,
	}, nil
}

func (f *fakeMessageService) MarkMessageAsRead(messageID string) (*domain.Message, error) {
	f.markCalledWithID = messageID
	if f.markErr != nil {
		return nil, f.markErr
	}
	return f.markReturnMsg, nil
}

func newTestContext(method, target, body string) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	var req *http.Request
	if reqBody != nil {
		req = httptest.NewRequest(method, target, reqBody)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestMessageHandler_CreateMessage_Success(t *testing.T) {
	fake := &fakeMessageService{}
	h := NewMessageHandler(fake)

	body := `{"senderId":"s1","recipientId":"r1","content":"hello"}`
	c, rec := newTestContext(http.MethodPost, "/messages", body)

	if err := h.CreateMessage(c); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusCreated)
	}

	var got domain.Message
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != "fake-id" {
		t.Errorf("ID = %q, want %q", got.ID, "fake-id")
	}
	if got.Content != "hello" {
		t.Errorf("Content = %q, want %q", got.Content, "hello")
	}
}

func TestMessageHandler_CreateMessage_IgnoresClientSuppliedServerFields(t *testing.T) {
	// The handler binds into CreateMessageRequest, which has no fields for
	// id/isRead/createdAt. Whatever the client sent for those keys should be
	// silently dropped at the binding step.
	fake := &fakeMessageService{}
	h := NewMessageHandler(fake)

	body := `{
		"id": "spoofed",
		"isRead": true,
		"createdAt": "2020-01-01T00:00:00Z",
		"senderId": "s1",
		"recipientId": "r1",
		"content": "hi"
	}`
	c, _ := newTestContext(http.MethodPost, "/messages", body)
	if err := h.CreateMessage(c); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if fake.createCalledWith == nil {
		t.Fatal("expected service to be called")
	}
	if fake.createCalledWith.SenderID != "s1" || fake.createCalledWith.Content != "hi" {
		t.Errorf("unexpected DTO: %+v", fake.createCalledWith)
	}
}

func TestMessageHandler_CreateMessage_BindError(t *testing.T) {
	fake := &fakeMessageService{}
	h := NewMessageHandler(fake)

	c, rec := newTestContext(http.MethodPost, "/messages", `{not valid json`)

	if err := h.CreateMessage(c); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestMessageHandler_CreateMessage_ServiceError(t *testing.T) {
	fake := &fakeMessageService{createErr: errors.New("db down")}
	h := NewMessageHandler(fake)

	body := `{"senderId":"s1","recipientId":"r1","content":"hello"}`
	c, rec := newTestContext(http.MethodPost, "/messages", body)

	if err := h.CreateMessage(c); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestMessageHandler_CreateMessage_ValidationError(t *testing.T) {
	fake := &fakeMessageService{createErr: fmt.Errorf("create message: %w", service.ErrValidation)}
	h := NewMessageHandler(fake)

	body := `{"senderId":"s1","recipientId":"r1","content":"hello"}`
	c, rec := newTestContext(http.MethodPost, "/messages", body)

	if err := h.CreateMessage(c); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	// Comma-ok map lookup: ok is true only if the key is actually present.
	// !ok means the "error" key was missing from the response body.
	if _, ok := got["error"]; !ok {
		t.Errorf("expected response body to have 'error' key, got %+v", got)
	}
}

func TestMessageHandler_MarkMessageAsRead_Success(t *testing.T) {
	fake := &fakeMessageService{
		markReturnMsg: &domain.Message{ID: "abc", IsRead: true, Content: "hi"},
	}
	h := NewMessageHandler(fake)

	c, rec := newTestContext(http.MethodPatch, "/messages/abc/read", "")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "abc"}})

	if err := h.MarkMessageAsRead(c); err != nil {
		t.Fatalf("MarkMessageAsRead returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	if fake.markCalledWithID != "abc" {
		t.Errorf("service called with id = %q, want %q", fake.markCalledWithID, "abc")
	}

	var got domain.Message
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !got.IsRead {
		t.Error("expected IsRead to be true")
	}
}

func TestMessageHandler_MarkMessageAsRead_ServiceError(t *testing.T) {
	fake := &fakeMessageService{markErr: errors.New("boom")}
	h := NewMessageHandler(fake)

	c, rec := newTestContext(http.MethodPatch, "/messages/missing/read", "")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "missing"}})

	if err := h.MarkMessageAsRead(c); err != nil {
		t.Fatalf("MarkMessageAsRead returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestMessageHandler_MarkMessageAsRead_ValidationError(t *testing.T) {
	fake := &fakeMessageService{markErr: fmt.Errorf("mark message as read: %w", service.ErrValidation)}
	h := NewMessageHandler(fake)

	c, rec := newTestContext(http.MethodPatch, "/messages/bad/read", "")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "bad"}})

	if err := h.MarkMessageAsRead(c); err != nil {
		t.Fatalf("MarkMessageAsRead returned error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestMessageHandler_MarkMessageAsRead_NotFound(t *testing.T) {
	fake := &fakeMessageService{markErr: fmt.Errorf("mark message as read: %w", service.ErrNotFound)}
	h := NewMessageHandler(fake)

	c, rec := newTestContext(http.MethodPatch, "/messages/missing/read", "")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "missing"}})

	if err := h.MarkMessageAsRead(c); err != nil {
		t.Fatalf("MarkMessageAsRead returned error: %v", err)
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
