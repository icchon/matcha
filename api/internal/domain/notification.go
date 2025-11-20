package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID          int64            `db:"id"`
	RecipientID uuid.UUID        `db:"recipient_id"`
	SenderID    sql.NullString   `db:"sender_id"` // ON DELETE SET NULL のため NullString (UUID)
	Type        NotificationType `db:"type"`
	IsRead      sql.NullBool     `db:"is_read"`
	CreatedAt   time.Time        `db:"created_at"`
}

type NotificationQueryRepository interface {
	Find(ctx context.Context, notificationID int64) (*Notification, error)
	GetByRecipientID(ctx context.Context, recipientID uuid.UUID) ([]Notification, error)
}

type NotificationCommandRepository interface {
	Save(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, notificationID int64) error
}

type NotificationRepository interface {
	NotificationQueryRepository
	NotificationCommandRepository
}
