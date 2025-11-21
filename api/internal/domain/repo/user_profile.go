package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserProfileQuery struct {
	UserID           *uuid.UUID
	FirstName        *string
	LastName         *string
	Username         *string
	Gender           *entity.Gender
	SexualPreference *entity.SexualPreference
	Biography        *string
	FameRating       *int32
	LocationName     *string
}

type CreateUserProfileParams struct {
	UserID           uuid.UUID
	FirstName        sql.NullString
	LastName         sql.NullString
	Username         sql.NullString
	Gender           sql.NullString
	SexualPreference sql.NullString
	Biography        sql.NullString
	LocationName     sql.NullString
}

type UserProfileQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Query(ctx context.Context, q *UserProfileQuery) ([]*entity.UserProfile, error)
}

type UserProfileCommandRepository interface {
	Create(ctx context.Context, params CreateUserProfileParams) (*entity.UserProfile, error)
	Update(ctx context.Context, userProfile *entity.UserProfile) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserProfileRepository interface {
	UserProfileQueryRepository
	UserProfileCommandRepository
}
