package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type PasswordReset struct {
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

type PasswordResetQueryRepository interface {
	Find(ctx context.Context, token string) (*PasswordReset, error)
}

type PasswordResetCommandRepository interface {
	Save(ctx context.Context, passwordReset *PasswordReset) error
	Delete(ctx context.Context, token string) error
}

type PasswordResetRepository interface {
	PasswordResetQueryRepository
	PasswordResetCommandRepository
}
