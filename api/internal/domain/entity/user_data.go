package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserData struct {
	UserID        uuid.UUID       `db:"user_id"`
	Latitude      sql.NullFloat64 `db:"latitude"`  // DECIMAL(10, 8)
	Longitude     sql.NullFloat64 `db:"longitude"` // DECIMAL(11, 8)
	InternalScore sql.NullInt32   `db:"internal_score"`
}
