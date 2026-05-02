package service

import (
	"errors"
	"strings"
	"testing"

	"hs-messaging-service/internal/domain"
)

// These tests exercise validation-only paths in MessageService. They rely on
// the fact that validation runs before the repository is touched, so a nil
// *postgres.MessageRepository never gets dereferenced. A future PR that
// introduces a repo interface would let us cover success paths here too.

func newMessageServiceForValidation() *MessageService {
	return &MessageService{messageRepository: nil}
}

func validUUID() string { return "11111111-1111-1111-1111-111111111111" }

func otherUUID() string { return "22222222-2222-2222-2222-222222222222" }

func TestValidateNewMessage_AllErrorsAreValidation(t *testing.T) {
	cases := []struct {
		name string
		msg  *domain.Message
	}{
		{"empty sender", &domain.Message{SenderID: "", RecipientID: validUUID(), Content: "hi"}},
		{"empty recipient", &domain.Message{SenderID: validUUID(), RecipientID: "", Content: "hi"}},
		{"empty content", &domain.Message{SenderID: validUUID(), RecipientID: otherUUID(), Content: ""}},
		{"invalid sender uuid", &domain.Message{SenderID: "not-a-uuid", RecipientID: validUUID(), Content: "hi"}},
		{"invalid recipient uuid", &domain.Message{SenderID: validUUID(), RecipientID: "not-a-uuid", Content: "hi"}},
		{"invalid job uuid", &domain.Message{SenderID: validUUID(), RecipientID: otherUUID(), Content: "hi", JobID: ptr("nope")}},
		{"self message", &domain.Message{SenderID: validUUID(), RecipientID: validUUID(), Content: "hi"}},
		{"content too long", &domain.Message{SenderID: validUUID(), RecipientID: otherUUID(), Content: strings.Repeat("a", maxContentLength+1)}},
	}

	svc := newMessageServiceForValidation()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.CreateMessage(tc.msg)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !errors.Is(err, ErrValidation) {
				t.Errorf("error %v should wrap ErrValidation", err)
			}
		})
	}
}

func TestValidateNewMessage_ContentAtMaxLengthIsAccepted(t *testing.T) {
	// We want to prove that content of exactly maxContentLength bytes passes
	// validation. The service is built with a nil repository (see
	// newMessageServiceForValidation), so once validation passes,
	// CreateMessage will call s.messageRepository.CreateMessage(...) on a nil
	// pointer, which crashes the goroutine via a runtime panic (Go's version
	// of an unrecoverable error / unhandled exception in other languages).
	//
	// In Go you can intercept a panic with recover(), but recover() only
	// works when called from a deferred function. `defer` schedules a
	// function to run when the surrounding function returns -- including when
	// it returns because of a panic. So the pattern below means:
	//
	//   1. Register a cleanup function that runs at the very end of the test.
	//   2. Inside it, call recover(): if a panic is in flight, it stops the
	//      panic and returns the panic value; if not, it returns nil.
	//   3. We discard the value with `_ =` because we don't care which kind
	//      of panic it was -- any panic here means we got past validation,
	//      which is exactly what we wanted to prove.
	//
	// Without this defer+recover, the nil-pointer dereference would mark the
	// test as FAILED even though validation behaved correctly.
	svc := newMessageServiceForValidation()
	msg := &domain.Message{
		SenderID:    validUUID(),
		RecipientID: otherUUID(),
		Content:     strings.Repeat("a", maxContentLength),
	}
	defer func() {
		_ = recover()
	}()
	err := svc.CreateMessage(msg)
	// We only reach this line if validation rejected the input before the
	// repo call. In that case, assert it wasn't an ErrValidation rejection.
	if err != nil && errors.Is(err, ErrValidation) {
		t.Errorf("content of exact max length should not fail validation, got %v", err)
	}
}

func TestMarkMessageAsRead_Validation(t *testing.T) {
	svc := newMessageServiceForValidation()

	t.Run("empty id", func(t *testing.T) {
		_, err := svc.MarkMessageAsRead("")
		if err == nil || !errors.Is(err, ErrValidation) {
			t.Errorf("expected validation error, got %v", err)
		}
	})

	t.Run("non-uuid id", func(t *testing.T) {
		_, err := svc.MarkMessageAsRead("not-a-uuid")
		if err == nil || !errors.Is(err, ErrValidation) {
			t.Errorf("expected validation error, got %v", err)
		}
	})
}

func ptr(s string) *string { return &s }
