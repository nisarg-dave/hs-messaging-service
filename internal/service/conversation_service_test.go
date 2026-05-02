package service

import (
	"errors"
	"testing"
)

func newConversationServiceForValidation() *ConversationService {
	return &ConversationService{conversationRepository: nil}
}

func TestListConversations_Validation(t *testing.T) {
	svc := newConversationServiceForValidation()

	cases := []struct {
		name   string
		userID string
	}{
		{"empty", ""},
		{"non-uuid", "not-a-uuid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.ListConversations(tc.userID)
			if err == nil || !errors.Is(err, ErrValidation) {
				t.Errorf("expected validation error, got %v", err)
			}
		})
	}
}

func TestGetConversation_Validation(t *testing.T) {
	svc := newConversationServiceForValidation()

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
			_, err := svc.GetConversation(tc.userID, tc.otherID)
			if err == nil || !errors.Is(err, ErrValidation) {
				t.Errorf("expected validation error, got %v", err)
			}
		})
	}
}
