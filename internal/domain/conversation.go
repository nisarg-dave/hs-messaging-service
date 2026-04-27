package domain

import "time"

type LastMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ConversationSummary struct {
	UserID      string      `json:"userId"`
	LastMessage LastMessage `json:"lastMessage"`
	UnreadCount int64       `json:"unreadCount"`
}
