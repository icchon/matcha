package domain

import (
	"context"
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

type MessageQueryRepository interface {
	Find(ctx context.Context, messageID int64) (*Message, error)
	GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]Message, error)
}

type MessageCommandRepository interface {
	Save(ctx context.Context, message *Message) error
	Delete(ctx context.Context, messageID int64) error
}

type MessageRepository interface {
	MessageQueryRepository
	MessageCommandRepository
}
