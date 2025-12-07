package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserProfile struct {
	UserID           uuid.UUID       `db:"user_id"`
	FirstName        sql.NullString  `db:"first_name"`
	LastName         sql.NullString  `db:"last_name"`
	Username         sql.NullString  `db:"username"`
	Gender           sql.NullString  `db:"gender"`
	SexualPreference sql.NullString  `db:"sexual_preference"`
	Birthday         sql.NullTime    `db:"birthday"`
	Occupation       sql.NullString  `db:"occupation"`
	Biography        sql.NullString  `db:"biography"`
	FameRating       sql.NullInt32   `db:"fame_rating"`
	LocationName     sql.NullString  `db:"location_name"`
	Distance         sql.NullFloat64 `db:"distance"`
}
