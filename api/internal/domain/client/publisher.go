package client

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessagePayload struct {
	ID          int64     `json:"id"`
	SenderID    uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
}

type NotificationPayload struct {
	ID          int64     `json:"id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	SenderID    string    `json:"sender_id"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

type AckPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	MessageID int64     `json:"message_id"`
	Timestamp int64     `json:"timestamp"`
}

type PresencePayload struct {
	UserID      uuid.UUID `json:"user_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Status      string    `json:"status"`
}

type ReadPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Timestamp   int64     `json:"timestamp"`
}

type Publisher interface {
	Publish(ctx context.Context, data interface{}) error
}
