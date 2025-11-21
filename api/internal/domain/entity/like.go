package entity

import (
	"github.com/google/uuid"
	"time"
)

type Like struct {
	LikerID   uuid.UUID `db:"liker_id"`
	LikedID   uuid.UUID `db:"liked_id"`
	CreatedAt time.Time `db:"created_at"`
}
