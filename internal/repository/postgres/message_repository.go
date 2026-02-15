package postgres

import (
	"log"

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

	log.Printf("Inserted %d rows", result.RowsAffected)
	return nil
}

func (r *MessageRepository) MarkMessageAsRead(messageID string) error {
	result := r.db.Model(&domain.Message{}).Where("id = ?", messageID).Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Marked message as read: %s", messageID)
	return nil
}