package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type AuthQuery struct {
	UserID      *uuid.UUID
	Email       *sql.NullString
	Provider    *entity.AuthProvider
	ProviderUID *sql.NullString
	IsVerified  *bool
}

type CreateAuthParams struct {
	UserID       uuid.UUID           `db:"user_id"`
	Email        sql.NullString      `db:"email"`
	Provider     entity.AuthProvider `db:"provider"`
	ProviderUID  sql.NullString      `db:"provider_uid"`
	IsVerified   bool                `db:"is_verified"`
	PasswordHash sql.NullString      `db:"password_hash"`
}

type AuthQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID, provider entity.AuthProvider) (*entity.Auth, error)
	Query(ctx context.Context, q *AuthQuery) ([]*entity.Auth, error)
}

type AuthCommandRepository interface {
	Create(ctx context.Context, auth *entity.Auth) error
	Update(ctx context.Context, auth *entity.Auth) error
	Delete(ctx context.Context, userID uuid.UUID, provider entity.AuthProvider) error
}

type AuthRepository interface {
	AuthQueryRepository
	AuthCommandRepository
}
