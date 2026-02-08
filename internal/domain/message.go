package domain

import "time"

type Message struct {
	ID        string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SenderID  string `json:"senderId" gorm:"type:uuid;not null;index"`
	RecipientID string `json:"recipientId" gorm:"type:uuid;not null;index"`
	Content   string `json:"content" gorm:"type:text;not null"`
	JobID     *string `json:"jobId,omitempty" gorm:"type:uuid;default:null"`
	IsRead    bool `json:"isRead" gorm:"type:boolean;default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"type:timestamp;default:now()"`
}