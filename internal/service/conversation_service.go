package service

import (
	"fmt"

	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/repository/postgres"
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
	return s.conversationRepository.ListConversations(userID)
}

func (s *ConversationService) GetConversation(userID, otherID string) ([]domain.Message, error) {
	if userID == "" {
		return nil, fmt.Errorf("get conversation: %w", errEmptyUserID)
	}
	if otherID == "" {
		return nil, fmt.Errorf("get conversation: %w", errEmptyOtherID)
	}
	return s.conversationRepository.GetConversation(userID, otherID)
}

var (
	errEmptyUserID  = fmt.Errorf("userID is required")
	errEmptyOtherID = fmt.Errorf("other userID is required")
)
