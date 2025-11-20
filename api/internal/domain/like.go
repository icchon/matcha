package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Like struct {
	LikerID   uuid.UUID `db:"liker_id"`
	LikedID   uuid.UUID `db:"liked_id"`
	CreatedAt time.Time `db:"created_at"`
}

type LikeQueryRepository interface {
	Find(ctx context.Context, likerID, likedID uuid.UUID) (*Like, error)
	GetByLikerID(ctx context.Context, likerID uuid.UUID) ([]Like, error)
}

type LikeCommandRepository interface {
	Save(ctx context.Context, like *Like) error
	Delete(ctx context.Context, likerID, likedID uuid.UUID) error
}

type LikeRepository interface {
	LikeQueryRepository
	LikeCommandRepository
}
