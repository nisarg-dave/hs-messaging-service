package service

import (
	"fmt"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/repository/postgres"

	"github.com/google/uuid"
)

type ConversationService struct {
	conversationRepository *postgres.ConversationRepository
}

func NewConversationService(conversationRepository *postgres.ConversationRepository) *ConversationService {
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
