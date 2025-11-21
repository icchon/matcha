package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type Auth struct {
	ID           int            `db:"id"`
	UserID       uuid.UUID      `db:"user_id"`
	Email        sql.NullString `db:"email"`
	Provider     AuthProvider   `db:"provider"`
	ProviderUID  sql.NullString `db:"provider_uid"`
	IsVerified   bool           `db:"is_verified"`
	PasswordHash sql.NullString `db:"password_hash"`
}
