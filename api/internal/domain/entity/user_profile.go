package entity

import (
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
