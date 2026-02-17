package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserData struct {
	UserID        uuid.UUID       `db:"user_id" json:"user_id"`
	Latitude      sql.NullFloat64 `db:"latitude" json:"latitude"`
	Longitude     sql.NullFloat64 `db:"longitude" json:"longitude"`
	InternalScore sql.NullInt32   `db:"internal_score" json:"internal_score"`
}
