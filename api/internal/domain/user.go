package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID    `db:"id"`
	CreatedAt      time.Time    `db:"created_at"`
	LastConnection sql.NullTime `db:"last_connection"`
}

type UserQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*User, error)
}

type UserCommandRepository interface {
	Save(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserRepository interface {
	UserQueryRepository
	UserCommandRepository
}
