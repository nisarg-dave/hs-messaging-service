package service

import (
	"fmt"

	"hs-messaging-service/internal/domain"

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
}

func NewConversationService(conversationRepository ConversationRepository) *ConversationService {
	return &ConversationService{conversationRepository: conversationRepository}
}

func (s *ConversationService) ListConversations(userID string) ([]domain.ConversationSummary, error) {
	if userID == "" {
		return nil, fmt.Errorf("list conversations: %w", errEmptyUserID)
	}
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("list conversations: %w", errInvalidUUID)
	}
	return s.conversationRepository.ListConversations(userID)
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
	return s.conversationRepository.GetConversation(userID, otherID)
}
