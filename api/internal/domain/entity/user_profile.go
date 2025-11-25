package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserProfile struct {
	UserID           uuid.UUID       `db:"user_id"`
	FirstName        sql.NullString  `db:"first_name"` // NULLを許容
	LastName         sql.NullString  `db:"last_name"`  // NULLを許容
	Username         sql.NullString  `db:"username"`   // NULLを許容
	Gender           sql.NullString  `db:"gender"`     // ENUM型
	SexualPreference sql.NullString  `db:"sexual_preference" json:"sexual_preference"`
	Birthday         sql.NullTime    `db:"birthday" json:"birthday"`
	Occupation       sql.NullString  `db:"occupation" json:"occupation"`
	Biography        sql.NullString  `db:"biography" json:"biography"`
	FameRating       sql.NullInt32   `db:"fame_rating" json:"fame_rating"`
	LocationName     sql.NullString  `db:"location_name" json:"location_name"`
	Distance         sql.NullFloat64 `db:"distance" json:"distance,omitempty"`
}
