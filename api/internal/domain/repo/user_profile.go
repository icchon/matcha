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
	UserID           uuid.UUID      `db:"user_id"`
	FirstName        sql.NullString `db:"first_name"`
	LastName         sql.NullString `db:"last_name"`
	Username         sql.NullString `db:"username"`
	Gender           sql.NullString `db:"gender"`
	SexualPreference sql.NullString `db:"sexual_preference"`
	Biography        sql.NullString `db:"biography"`
	LocationName     sql.NullString `db:"location_name"`
}

type UserProfileQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Query(ctx context.Context, q *UserProfileQuery) ([]*entity.UserProfile, error)
}

type UserProfileCommandRepository interface {
	Create(ctx context.Context, userProfile *entity.UserProfile) (error)
	Update(ctx context.Context, userProfile *entity.UserProfile) (error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserProfileRepository interface {
	UserProfileQueryRepository
	UserProfileCommandRepository
}
