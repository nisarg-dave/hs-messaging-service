package service

import (
	"errors"
	"testing"

	"hs-messaging-service/internal/domain"
)

type fakeConversationRepository struct {
	listCalledWith string
	listReturn     []domain.ConversationSummary
	listErr        error

	getCalledWithUser  string
	getCalledWithOther string
	getReturn          []domain.Message
	getErr             error
}

func (f *fakeConversationRepository) ListConversations(userID string) ([]domain.ConversationSummary, error) {
	f.listCalledWith = userID
	return f.listReturn, f.listErr
}

func (f *fakeConversationRepository) GetConversation(userID, otherID string) ([]domain.Message, error) {
	f.getCalledWithUser = userID
	f.getCalledWithOther = otherID
	return f.getReturn, f.getErr
}

func TestListConversations_Success(t *testing.T) {
	want := []domain.ConversationSummary{{UserID: otherUUID(), UnreadCount: 3}}
	repo := &fakeConversationRepository{listReturn: want}
	svc := NewConversationService(repo, testLogger())

	got, err := svc.ListConversations(validUUID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].UserID != otherUUID() {
		t.Errorf("unexpected result: %+v", got)
	}
	if repo.listCalledWith != validUUID() {
		t.Errorf("repo called with %q, want %q", repo.listCalledWith, validUUID())
	}
}

func TestListConversations_Validation(t *testing.T) {
	cases := []struct {
		name   string
		userID string
	}{
		{"empty", ""},
		{"non-uuid", "not-a-uuid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeConversationRepository{}
			svc := NewConversationService(repo, testLogger())
			_, err := svc.ListConversations(tc.userID)
			if err == nil || !errors.Is(err, ErrValidation) {
				t.Errorf("expected validation error, got %v", err)
			}
			if repo.listCalledWith != "" {
				t.Error("validation should short-circuit before the repo is called")
			}
		})
	}
}

func TestGetConversation_Success(t *testing.T) {
	want := []domain.Message{{ID: "m1", Content: "hi"}}
	repo := &fakeConversationRepository{getReturn: want}
	svc := NewConversationService(repo, testLogger())

	got, err := svc.GetConversation(validUUID(), otherUUID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "m1" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestGetConversation_Validation(t *testing.T) {
	cases := []struct {
		name            string
		userID, otherID string
	}{
		{"empty user", "", otherUUID()},
		{"empty other", validUUID(), ""},
		{"non-uuid user", "nope", otherUUID()},
		{"non-uuid other", validUUID(), "nope"},
		{"self conversation", validUUID(), validUUID()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeConversationRepository{}
			svc := NewConversationService(repo, testLogger())
			_, err := svc.GetConversation(tc.userID, tc.otherID)
			if err == nil || !errors.Is(err, ErrValidation) {
				t.Errorf("expected validation error, got %v", err)
			}
			if repo.getCalledWithUser != "" {
				t.Error("validation should short-circuit before the repo is called")
			}
		})
	}
}
