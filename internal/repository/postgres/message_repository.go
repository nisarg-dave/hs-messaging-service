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

func (r *MessageRepository) MarkMessageAsRead(messageID string) (*domain.Message, error) {
	result := r.db.Model(&domain.Message{}).Where("id = ?", messageID).Update("is_read", true)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Printf("Marked message as read: %s", messageID)

	// new(T) allocates a zero-value T on the heap and returns *T.
	// So new(domain.Message) is equivalent to &domain.Message{} — both give
	// a pointer to an empty Message that GORM's First can fill in.
	// First takes a destination pointer; message is already *Message, so
	// &message is **Message — GORM accepts that and writes into *message.
	message := new(domain.Message)
	result = r.db.First(&message, "id = ?", messageID)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}
