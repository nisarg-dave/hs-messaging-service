package service

import (
	"fmt"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/repository/postgres"

	"github.com/google/uuid"
)

type MessageService struct {
	messageRepository *postgres.MessageRepository
}

func NewMessageService(messageRepository *postgres.MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (s *MessageService) CreateMessage(message *domain.Message) error {
	if err := validateNewMessage(message); err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return s.messageRepository.CreateMessage(message)
}

func (s *MessageService) MarkMessageAsRead(messageID string) (*domain.Message, error) {
	if messageID == "" {
		return nil, fmt.Errorf("mark message as read: %w", errEmptyMessageID)
	}
	if _, err := uuid.Parse(messageID); err != nil {
		return nil, fmt.Errorf("mark message as read: %w", errInvalidUUID)
	}
	message, err := s.messageRepository.MarkMessageAsRead(messageID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func validateNewMessage(m *domain.Message) error {
	if m == nil {
		return errEmptyContent
	}

	_, senderErr := uuid.Parse(m.SenderID)
	_, recipientErr := uuid.Parse(m.RecipientID)

	switch {
	case m.SenderID == "":
		return errEmptySenderID
	case m.RecipientID == "":
		return errEmptyRecipientID
	case m.Content == "":
		return errEmptyContent
	case senderErr != nil, recipientErr != nil:
		return errInvalidUUID
	case m.JobID != nil && !isUUID(*m.JobID):
		return errInvalidUUID
	case m.SenderID == m.RecipientID:
		return errSelfMessage
	case len(m.Content) > maxContentLength:
		return errContentTooLong
	}
	return nil
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
