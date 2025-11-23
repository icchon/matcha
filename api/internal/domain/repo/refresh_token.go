package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type RefreshTokenQuery struct {
	TokenHash *string
	UserID    *uuid.UUID
	Revoked   *bool
}

type RefreshTokenQueryRepository interface {
	Find(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Query(ctx context.Context, q *RefreshTokenQuery) ([]*entity.RefreshToken, error)
}

type RefreshTokenCommandRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	Update(ctx context.Context, token *entity.RefreshToken) error
	Delete(ctx context.Context, tokenHash string) error
}

type RefreshTokenRepository interface {
	RefreshTokenQueryRepository
	RefreshTokenCommandRepository
}
