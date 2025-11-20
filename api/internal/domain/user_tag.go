package domain

import (
	"context"
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

type UserTagQueryRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Tag, error)
}

type UserTagCommandRepository interface {
	Save(ctx context.Context, userTag *UserTag) error
	Delete(ctx context.Context, userID uuid.UUID, tagID int32) error
}

type UserTagRepository interface {
	UserTagQueryRepository
	UserTagCommandRepository
}
