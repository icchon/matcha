package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID    `db:"id"`
	CreatedAt      time.Time    `db:"created_at"`
	LastConnection sql.NullTime `db:"last_connection"`
}
