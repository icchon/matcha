package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type UserProfile struct {
	UserID           uuid.UUID      `db:"user_id"`
	FirstName        sql.NullString `db:"first_name"`        // NULLを許容
	LastName         sql.NullString `db:"last_name"`         // NULLを許容
	Username         sql.NullString `db:"username"`          // NULLを許容
	Gender           sql.NullString `db:"gender"`            // ENUM型
	SexualPreference sql.NullString `db:"sexual_preference"` // ENUM型
	Biography        sql.NullString `db:"biography"`
	FameRating       sql.NullInt32  `db:"fame_rating"`
	LocationName     sql.NullString `db:"location_name"`
}

type UserProfileQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
}

type UserProfileCommandRepository interface {
	Save(ctx context.Context, userProfile *UserProfile) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserProfileRepository interface {
	UserProfileQueryRepository
	UserProfileCommandRepository
}
