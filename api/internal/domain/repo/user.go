package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserQuery struct {
	ID *uuid.UUID
}

type UserQueryRepository interface {
	Find(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	Query(ctx context.Context, q *UserQuery) ([]*entity.User, error)
}

type UserCommandRepository interface {
	Create(ctx context.Context, user* entity.User) (error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserRepository interface {
	UserQueryRepository
	UserCommandRepository
}
