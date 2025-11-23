package entity

import (
	"github.com/google/uuid"
	"time"
)

type Connection struct {
	User1ID   uuid.UUID `db:"user1_id" json:"user1_id"`
	User2ID   uuid.UUID `db:"user2_id" json:"user2_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
