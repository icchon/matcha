package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type LikeQuery struct {
	LikerID   *uuid.UUID
	LikedID   *uuid.UUID
	CreatedAt *time.Time
}

type LikeQueryRepository interface {
	Find(ctx context.Context, likerID, likedID uuid.UUID) (*entity.Like, error)
	Query(ctx context.Context, q *LikeQuery) ([]*entity.Like, error)
}

type LikeCommandRepository interface {
	Create(ctx context.Context, like *entity.Like) error
	Delete(ctx context.Context, likerID, likedID uuid.UUID) error
}

type LikeRepository interface {
	LikeQueryRepository
	LikeCommandRepository
}
