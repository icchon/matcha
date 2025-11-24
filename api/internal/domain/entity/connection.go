package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
)

type Connection struct {
	User1ID   uuid.UUID `db:"user1_id" json:"user1_id"`
	User2ID   uuid.UUID `db:"user2_id" json:"user2_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (c *Connection) GetOtherUserID(user1D uuid.UUID) (uuid.UUID, error) {
	if c.User1ID == user1D {
		return c.User2ID, nil
	}
	if c.User2ID == user1D {
		return c.User1ID, nil
	}
	return uuid.Nil, apperrors.ErrNotFound
}
