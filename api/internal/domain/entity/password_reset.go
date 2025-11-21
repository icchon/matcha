package entity

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
