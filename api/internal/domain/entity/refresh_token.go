package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	TokenHash string    `db:"token_hash"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	Revoked   bool      `db:"revoked"`
}
