package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserTagQuery struct {
	UserID  *uuid.UUID
	TagID   *int32
	TagName *string
}

type UserTagQueryRepository interface {
	Query(ctx context.Context, q *UserTagQuery) ([]*entity.Tag, error)
}

type UserTagCommandRepository interface {
	Create(ctx context.Context, userTag *entity.UserTag) error
	Delete(ctx context.Context, userID uuid.UUID, tagID int32) error
}

type UserTagRepository interface {
	UserTagQueryRepository
	UserTagCommandRepository
}
