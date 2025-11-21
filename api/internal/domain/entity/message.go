package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID          int64        `db:"id"`
	SenderID    uuid.UUID    `db:"sender_id"`
	RecipientID uuid.UUID    `db:"recipient_id"`
	Content     string       `db:"content"`
	SentAt      time.Time    `db:"sent_at"`
	IsRead      sql.NullBool `db:"is_read"`
}
