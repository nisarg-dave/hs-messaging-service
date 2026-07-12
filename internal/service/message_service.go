package service

import (
	"errors"
	"fmt"

	"hs-messaging-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageRepository is the slice of repository behavior MessageService needs.
// Defining it here (the consumer) lets tests inject fakes without dragging in
// a real Postgres connection. *postgres.MessageRepository satisfies this
// interface implicitly because it has matching method signatures.
type MessageRepository interface {
	CreateMessage(message *domain.Message) error
	MarkMessageAsRead(messageID string) (*domain.Message, error)
}

// CreateMessageRequest is the input shape the service accepts for creating a
// new message. It deliberately omits server-controlled fields (ID, IsRead,
// CreatedAt, UpdatedAt) so a malicious or buggy client can't spoof them. The
// handler binds JSON into this struct and the service builds a domain.Message
// from it before persisting.
type CreateMessageRequest struct {
	SenderID    string  `json:"senderId"`
	RecipientID string  `json:"recipientId"`
	Content     string  `json:"content"`
	JobID       *string `json:"jobId,omitempty"`
}

type MessageService struct {
	messageRepository MessageRepository
}

func NewMessageService(messageRepository MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

// CreateMessage validates the request and, on success, persists a new message.
// Returns the persisted domain.Message so callers can see the server-assigned
// ID and timestamps.
func (s *MessageService) CreateMessage(req *CreateMessageRequest) (*domain.Message, error) {
	// (s *MessageService) above is a method receiver: this is a method on
	// *MessageService, so you call it as svc.CreateMessage(...). The receiver s
	// is the instance the method runs on (like "this" in other languages).
	if err := validateCreateMessageRequest(req); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	message := &domain.Message{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
		Content:     req.Content,
		JobID:       req.JobID,
	}
	if err := s.messageRepository.CreateMessage(message); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	return message, nil
}

func (s *MessageService) MarkMessageAsRead(messageID string) (*domain.Message, error) {
	if messageID == "" {
		return nil, fmt.Errorf("mark message as read: %w", errEmptyMessageID)
	}
	if !isUUID(messageID) {
		return nil, fmt.Errorf("mark message as read: %w", errInvalidUUID)
	}
	message, err := s.messageRepository.MarkMessageAsRead(messageID)
	if err != nil {
		// Translate the GORM-specific not-found into our domain sentinel so
		// the handler can map it to 404 without importing gorm.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("mark message as read: %w", errMessageNotFound)
		}
		return nil, fmt.Errorf("mark message as read: %w", err)
	}
	return message, nil
}

func validateCreateMessageRequest(r *CreateMessageRequest) error {
	if r == nil {
		return errEmptyContent
	}

	// Empty checks first, then isUUID for format — same helper JobID already
	// uses. Parsing into discarded errors beforehand was redundant.
	switch {
	case r.SenderID == "":
		return errEmptySenderID
	case r.RecipientID == "":
		return errEmptyRecipientID
	case r.Content == "":
		return errEmptyContent
	case !isUUID(r.SenderID), !isUUID(r.RecipientID):
		return errInvalidUUID
	case r.JobID != nil && !isUUID(*r.JobID):
		return errInvalidUUID
	case r.SenderID == r.RecipientID:
		return errSelfMessage
	case len(r.Content) > maxContentLength:
		return errContentTooLong
	}
	return nil
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
