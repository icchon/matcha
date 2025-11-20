package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type VerificationToken struct {
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}

type VerificationTokenQueryRepository interface {
	Find(ctx context.Context, token string) (*VerificationToken, error)
}

type VerificationTokenCommandRepository interface {
	Save(ctx context.Context, token *VerificationToken) error
	Delete(ctx context.Context, token string) error
}

type VerificationTokenRepository interface {
	VerificationTokenQueryRepository
	VerificationTokenCommandRepository
}
