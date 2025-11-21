package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type PasswordResetQuery struct {
	UserID    *uuid.UUID
	Token     *string
	ExpiresAt *time.Time
}

type CreatePasswordResetParams struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type PasswordResetQueryRepository interface {
	Find(ctx context.Context, token string) (*entity.PasswordReset, error)
	Query(ctx context.Context, q *PasswordResetQuery) ([]*entity.PasswordReset, error)
}

type PasswordResetCommandRepository interface {
	Create(ctx context.Context, params CreatePasswordResetParams) (*entity.PasswordReset, error)
	Update(ctx context.Context, passwordReset *entity.PasswordReset) error
	Delete(ctx context.Context, token string) error
}

type PasswordResetRepository interface {
	PasswordResetQueryRepository
	PasswordResetCommandRepository
}
