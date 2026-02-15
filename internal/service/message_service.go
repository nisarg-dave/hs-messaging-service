package service

import (
	"hs-messaging-service/internal/domain"
	"hs-messaging-service/internal/repository/postgres"
)	

type MessageService struct {
	messageRepository *postgres.MessageRepository
}

func NewMessageService(messageRepository *postgres.MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (s *MessageService) CreateMessage(message *domain.Message) error {
	// Here we can do other things like validate the message, check for profanity, call other services, etc.
	return s.messageRepository.CreateMessage(message)
}

func (s *MessageService) MarkMessageAsRead(messageID string) error {
	return s.messageRepository.MarkMessageAsRead(messageID)
}