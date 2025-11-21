package entity

import (
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
