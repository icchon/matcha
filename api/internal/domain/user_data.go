package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type UserData struct {
	UserID        uuid.UUID       `db:"user_id"`
	Latitude      sql.NullFloat64 `db:"latitude"`  // DECIMAL(10, 8)
	Longitude     sql.NullFloat64 `db:"longitude"` // DECIMAL(11, 8)
	InternalScore sql.NullInt32   `db:"internal_score"`
}

type UserDataQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*UserData, error)
}

type UserDataCommandRepository interface {
	Save(ctx context.Context, userData *UserData) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserDataRepository interface {
	UserDataQueryRepository
	UserDataCommandRepository
}
