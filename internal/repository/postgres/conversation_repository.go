package postgres

import (
	"time"

	"hs-messaging-service/internal/domain"

	"gorm.io/gorm"
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

type conversationRow struct {
	UserID      string
	Content     string
	CreatedAt   time.Time
	UnreadCount int64
}

func (r *ConversationRepository) ListConversations(userID string) ([]domain.ConversationSummary, error) {
	const query = `
WITH paired AS (
    SELECT
        CASE WHEN sender_id = ? THEN recipient_id ELSE sender_id END AS other_user,
        content,
        created_at,
        recipient_id,
        is_read
    FROM messages
    WHERE sender_id = ? OR recipient_id = ?
),
ranked AS (
    SELECT
        other_user,
        content,
        created_at,
        ROW_NUMBER() OVER (PARTITION BY other_user ORDER BY created_at DESC) AS rn
    FROM paired
)
SELECT
    r.other_user AS user_id,
    r.content AS content,
    r.created_at AS created_at,
    COALESCE((
        SELECT COUNT(*) FROM paired p
        WHERE p.other_user = r.other_user
          AND p.recipient_id = ?
          AND p.is_read = false
    ), 0) AS unread_count
FROM ranked r
WHERE r.rn = 1
ORDER BY r.created_at DESC
`

	var rows []conversationRow
	if err := r.db.Raw(query, userID, userID, userID, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]domain.ConversationSummary, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.ConversationSummary{
			UserID: row.UserID,
			LastMessage: domain.LastMessage{
				Content:   row.Content,
				CreatedAt: row.CreatedAt,
			},
			UnreadCount: row.UnreadCount,
		})
	}
	return out, nil
}

func (r *ConversationRepository) GetConversation(userID, otherID string) ([]domain.Message, error) {
	var messages []domain.Message
	result := r.db.
		Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			userID, otherID, otherID, userID).
		Order("created_at ASC").
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
