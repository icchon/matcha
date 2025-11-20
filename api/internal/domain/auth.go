package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Auth struct {
	UserID       uuid.UUID      `db:"user_id"`
	Email        sql.NullString `db:"email"`
	Provider     AuthProvider   `db:"provider"`
	ProviderUID  string         `db:"provider_uid"`
	PasswordHash sql.NullString `db:"password_hash"` // 外部認証の場合はNULL
}

type AuthQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID, provider AuthProvider) (*Auth, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Auth, error)
}

type AuthCommandRepository interface {
	Save(ctx context.Context, auth *Auth) error
	Delete(ctx context.Context, userID uuid.UUID, provider AuthProvider) error
}

type AuthRepository interface {
	AuthQueryRepository
	AuthCommandRepository
}
