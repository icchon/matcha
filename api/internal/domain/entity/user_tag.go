package entity

import (
	"github.com/google/uuid"
)

type Tag struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type UserTag struct {
	UserID uuid.UUID `db:"user_id"`
	TagID  int32     `db:"tag_id"`
}
