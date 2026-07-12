package service

import (
	"fmt"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/logging"

	"github.com/google/uuid"
)

// ConversationRepository is the slice of repository behavior
// ConversationService needs. *postgres.ConversationRepository satisfies it
// implicitly via matching method signatures.
type ConversationRepository interface {
	ListConversations(userID string) ([]domain.ConversationSummary, error)
	GetConversation(userID, otherID string) ([]domain.Message, error)
}

type ConversationService struct {
	conversationRepository ConversationRepository
	logger                 logging.Logger
}

// NewConversationService wires dependencies via constructor injection — see
// NewMessageService for the SOLID/pattern rationale.
func NewConversationService(conversationRepository ConversationRepository, logger logging.Logger) *ConversationService {
	return &ConversationService{conversationRepository: conversationRepository, logger: logger}
}

func (s *ConversationService) ListConversations(userID string) ([]domain.ConversationSummary, error) {
	if userID == "" {
		return nil, fmt.Errorf("list conversations: %w", errEmptyUserID)
	}
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("list conversations: %w", errInvalidUUID)
	}
	summaries, err := s.conversationRepository.ListConversations(userID)
	if err != nil {
		return nil, err
	}
	s.logger.Info("conversations listed", "userId", userID, "count", len(summaries))
	return summaries, nil
}

func (s *ConversationService) GetConversation(userID, otherID string) ([]domain.Message, error) {
	if userID == "" {
		return nil, fmt.Errorf("get conversation: %w", errEmptyUserID)
	}
	if otherID == "" {
		return nil, fmt.Errorf("get conversation: %w", errEmptyOtherID)
	}
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("get conversation: %w", errInvalidUUID)
	}
	if _, err := uuid.Parse(otherID); err != nil {
		return nil, fmt.Errorf("get conversation: %w", errInvalidUUID)
	}
	if userID == otherID {
		return nil, fmt.Errorf("get conversation: %w", errSelfConversation)
	}
	messages, err := s.conversationRepository.GetConversation(userID, otherID)
	if err != nil {
		return nil, err
	}
	s.logger.Info("conversation retrieved",
		"userId", userID,
		"otherId", otherID,
		"messageCount", len(messages),
	)
	return messages, nil
}
