package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserDataQuery struct {
	UserID        *uuid.UUID
	Latitude      *float64
	Longitude     *float64
	InternalScore *int32
}

type UserDataQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*entity.UserData, error)
	Query(ctx context.Context, q *UserDataQuery) ([]*entity.UserData, error)
}

type UserDataCommandRepository interface {
	Create(ctx context.Context, userData *entity.UserData) (error)
	Update(ctx context.Context, userData *entity.UserData) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserDataRepository interface {
	UserDataQueryRepository
	UserDataCommandRepository
}
