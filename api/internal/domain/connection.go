package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Connection struct {
	User1ID   uuid.UUID `db:"user1_id"`
	User2ID   uuid.UUID `db:"user2_id"`
	CreatedAt time.Time `db:"created_at"`
}

type ConnectionQueryRepository interface {
	Find(ctx context.Context, user1ID, user2ID uuid.UUID) (*Connection, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Connection, error)
}

type ConnectionCommandRepository interface {
	Save(ctx context.Context, connection *Connection) error
	Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error
}

type ConnectionRepository interface {
	ConnectionQueryRepository
	ConnectionCommandRepository
}
