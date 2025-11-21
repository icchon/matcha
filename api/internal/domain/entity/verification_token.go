package entity

import (
	"github.com/google/uuid"
	"time"
)

type VerificationToken struct {
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
