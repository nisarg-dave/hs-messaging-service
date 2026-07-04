package postgres

import (
	"hs-messaging-service/internal/domain"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(message *domain.Message) error {
	result := r.db.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MessageRepository) MarkMessageAsRead(messageID string) (*domain.Message, error) {
	result := r.db.Model(&domain.Message{}).Where("id = ?", messageID).Update("is_read", true)
	if result.Error != nil {
		return nil, result.Error
	}

	message := new(domain.Message)
	result = r.db.First(&message, "id = ?", messageID)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}
