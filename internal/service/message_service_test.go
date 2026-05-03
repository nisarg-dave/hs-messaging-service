package service

import (
	"errors"
	"strings"
	"testing"

	"hs-messaging-service/internal/domain"

	"gorm.io/gorm"
)

// fakeMessageRepository implements the MessageRepository interface so we can
// drive MessageService in unit tests without a real database.
type fakeMessageRepository struct {
	createCalledWith *domain.Message
	createErr        error

	markCalledWith string
	markReturn     *domain.Message
	markErr        error
}

func (f *fakeMessageRepository) CreateMessage(m *domain.Message) error {
	// Snapshot a copy so tests can assert what the service passed in BEFORE
	// the repo mutates the struct (the real Postgres repo fills in ID on
	// insert, mirrored here).
	snapshot := *m
	f.createCalledWith = &snapshot
	if f.createErr != nil {
		return f.createErr
	}
	m.ID = "generated-id"
	return nil
}

func (f *fakeMessageRepository) MarkMessageAsRead(id string) (*domain.Message, error) {
	f.markCalledWith = id
	if f.markErr != nil {
		return nil, f.markErr
	}
	return f.markReturn, nil
}

func validUUID() string { return "11111111-1111-1111-1111-111111111111" }
func otherUUID() string { return "22222222-2222-2222-2222-222222222222" }

func validRequest() *CreateMessageRequest {
	return &CreateMessageRequest{
		SenderID:    validUUID(),
		RecipientID: otherUUID(),
		Content:     "hi",
	}
}

func TestCreateMessage_Success(t *testing.T) {
	repo := &fakeMessageRepository{}
	svc := NewMessageService(repo)

	msg, err := svc.CreateMessage(validRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID != "generated-id" {
		t.Errorf("expected repo to fill ID, got %q", msg.ID)
	}
	if repo.createCalledWith == nil || repo.createCalledWith.Content != "hi" {
		t.Errorf("repo not called with expected message: %+v", repo.createCalledWith)
	}
}

func TestCreateMessage_StripsServerControlledFields(t *testing.T) {
	// Even if a future caller sneaks server-controlled fields onto the DTO,
	// the service should never propagate them. The DTO doesn't expose them
	// at all, so this test mainly documents that the persisted domain.Message
	// starts from a zero ID/IsRead/CreatedAt — letting the DB defaults win.
	repo := &fakeMessageRepository{}
	svc := NewMessageService(repo)

	if _, err := svc.CreateMessage(validRequest()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createCalledWith.ID != "" {
		t.Errorf("expected zero ID before insert, got %q", repo.createCalledWith.ID)
	}
	if repo.createCalledWith.IsRead {
		t.Error("expected IsRead=false before insert")
	}
	if !repo.createCalledWith.CreatedAt.IsZero() {
		t.Errorf("expected zero CreatedAt before insert, got %v", repo.createCalledWith.CreatedAt)
	}
}

func TestCreateMessage_ContentAtMaxLengthIsAccepted(t *testing.T) {
	repo := &fakeMessageRepository{}
	svc := NewMessageService(repo)

	req := validRequest()
	req.Content = strings.Repeat("a", maxContentLength)

	if _, err := svc.CreateMessage(req); err != nil {
		t.Fatalf("expected content at exact max length to be accepted, got %v", err)
	}
}

func TestCreateMessage_ValidationFailures(t *testing.T) {
	cases := []struct {
		name string
		req  *CreateMessageRequest
	}{
		{"empty sender", &CreateMessageRequest{SenderID: "", RecipientID: validUUID(), Content: "hi"}},
		{"empty recipient", &CreateMessageRequest{SenderID: validUUID(), RecipientID: "", Content: "hi"}},
		{"empty content", &CreateMessageRequest{SenderID: validUUID(), RecipientID: otherUUID(), Content: ""}},
		{"invalid sender uuid", &CreateMessageRequest{SenderID: "not-a-uuid", RecipientID: validUUID(), Content: "hi"}},
		{"invalid recipient uuid", &CreateMessageRequest{SenderID: validUUID(), RecipientID: "not-a-uuid", Content: "hi"}},
		{"invalid job uuid", &CreateMessageRequest{SenderID: validUUID(), RecipientID: otherUUID(), Content: "hi", JobID: ptr("nope")}},
		{"self message", &CreateMessageRequest{SenderID: validUUID(), RecipientID: validUUID(), Content: "hi"}},
		{"content too long", &CreateMessageRequest{SenderID: validUUID(), RecipientID: otherUUID(), Content: strings.Repeat("a", maxContentLength+1)}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeMessageRepository{}
			svc := NewMessageService(repo)

			_, err := svc.CreateMessage(tc.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !errors.Is(err, ErrValidation) {
				t.Errorf("error %v should wrap ErrValidation", err)
			}
			if repo.createCalledWith != nil {
				t.Error("validation should short-circuit before the repo is called")
			}
		})
	}
}

func TestCreateMessage_RepoErrorPropagates(t *testing.T) {
	repo := &fakeMessageRepository{createErr: errors.New("db down")}
	svc := NewMessageService(repo)

	_, err := svc.CreateMessage(validRequest())
	if err == nil {
		t.Fatal("expected error from repo")
	}
	if errors.Is(err, ErrValidation) {
		t.Errorf("repo error should not be classified as validation: %v", err)
	}
}

func TestMarkMessageAsRead_Success(t *testing.T) {
	want := &domain.Message{ID: validUUID(), IsRead: true}
	repo := &fakeMessageRepository{markReturn: want}
	svc := NewMessageService(repo)

	got, err := svc.MarkMessageAsRead(validUUID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestMarkMessageAsRead_Validation(t *testing.T) {
	cases := []struct {
		name string
		id   string
	}{
		{"empty id", ""},
		{"non-uuid id", "not-a-uuid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewMessageService(&fakeMessageRepository{})
			_, err := svc.MarkMessageAsRead(tc.id)
			if err == nil || !errors.Is(err, ErrValidation) {
				t.Errorf("expected validation error, got %v", err)
			}
		})
	}
}

func TestMarkMessageAsRead_NotFound(t *testing.T) {
	repo := &fakeMessageRepository{markErr: gorm.ErrRecordNotFound}
	svc := NewMessageService(repo)

	_, err := svc.MarkMessageAsRead(validUUID())
	if err == nil || !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMarkMessageAsRead_OtherRepoErrorPropagates(t *testing.T) {
	repo := &fakeMessageRepository{markErr: errors.New("db down")}
	svc := NewMessageService(repo)

	_, err := svc.MarkMessageAsRead(validUUID())
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrValidation) {
		t.Errorf("generic repo error should not be classified: %v", err)
	}
}

func ptr(s string) *string { return &s }
