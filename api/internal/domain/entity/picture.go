package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Picture struct {
	ID           int32        `db:"id"`
	UserID       uuid.UUID    `db:"user_id"`
	URL          string       `db:"url"`
	IsProfilePic sql.NullBool `db:"is_profile_pic"`
	CreatedAt    time.Time    `db:"created_at"`
}
