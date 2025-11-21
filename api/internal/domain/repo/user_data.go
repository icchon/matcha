package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserDataQuery struct {
	UserID        *uuid.UUID
	Latitude      *float64
	Longitude     *float64
	InternalScore *int32
}

type CreateUserDataParams struct {
	UserID        uuid.UUID
	Latitude      sql.NullFloat64
	Longitude     sql.NullFloat64
	InternalScore sql.NullInt32
}

type UserDataQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*entity.UserData, error)
	Query(ctx context.Context, q *UserDataQuery) ([]*entity.UserData, error)
}

type UserDataCommandRepository interface {
	Create(ctx context.Context, params CreateUserDataParams) (*entity.UserData, error)
	Update(ctx context.Context, userData *entity.UserData) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserDataRepository interface {
	UserDataQueryRepository
	UserDataCommandRepository
}
