package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type VerificationTokenQuery struct {
	Token     *string
	UserID    *uuid.UUID
	ExpiresAt *time.Time
}

type VerificationTokenQueryRepository interface {
	Find(ctx context.Context, token string) (*entity.VerificationToken, error)
	Query(ctx context.Context, q *VerificationTokenQuery) ([]*entity.VerificationToken, error)
}

type VerificationTokenCommandRepository interface {
	Create(ctx context.Context, token *entity.VerificationToken) (error)
	Update(ctx context.Context, token *entity.VerificationToken) error
	Delete(ctx context.Context, token string) error
}

type VerificationTokenRepository interface {
	VerificationTokenQueryRepository
	VerificationTokenCommandRepository
}
