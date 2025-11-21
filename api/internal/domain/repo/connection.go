package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type ConnectionQuery struct {
	User1ID   *uuid.UUID
	User2ID   *uuid.UUID
	CreatedAt *time.Time
}

type ConnectionQueryRepository interface {
	Find(ctx context.Context, user1ID, user2ID uuid.UUID) (*entity.Connection, error)
	Query(ctx context.Context, q *ConnectionQuery) ([]*entity.Connection, error)
}

type ConnectionCommandRepository interface {
	Create(ctx context.Context, connection *entity.Connection) error
	Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error
}

type ConnectionRepository interface {
	ConnectionQueryRepository
	ConnectionCommandRepository
}
