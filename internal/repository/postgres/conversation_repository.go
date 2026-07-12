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
	// Conversation list SQL uses CTEs (WITH …): named chunks run top-to-bottom.
	//
	// paired: every message involving this user (sender OR recipient); CASE maps each row’s
	// peer to other_user so both directions share one thread key; keeps recipient_id and is_read
	// for the unread tally.
	//
	// ranked: ROW_NUMBER per other_user ordered by newest first → rn == 1 is the latest message
	// per 1:1 thread.
	//
	// Outer query: keeps only rn=1 (last preview) and attaches unread_count by counting rows in
	// paired where same peer, recipient_id is you, and is_read is false.
	//
	// userID is passed four times below because ? placeholders are positional: same id fills
	// CASE (sender check), WHERE sender, WHERE recipient, and unread recipient filter.
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
	// Raw executes the string SQL; each userID aligns with a ? in order. Scan fills rows —
	// one conversationRow per result row matching the SELECT aliases (user_id, content, …).
	if err := r.db.Raw(query, userID, userID, userID, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// make([]T, length, capacity) builds a slice. Here length is 0 (nothing
	// in it yet) and capacity is len(rows) so append won't need to grow the
	// underlying array as we loop — one allocation for the whole result.
	// Equivalent idea in JS: new Array(n) then push, but Go keeps length and
	// capacity separate.
	out := make([]domain.ConversationSummary, 0, len(rows))
	for _, row := range rows {
		// Nested struct literals: map flat SQL row → domain.ConversationSummary with nested
		// LastMessage (Content + CreatedAt together).
		// append lengthens out by one — always assign: out = append(out, elt). Similar in spirit to
		// [...prev, elt] spread in JS, but Go returns a slice you must capture.
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
