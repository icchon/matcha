package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type MessageQuery struct {
	ID          *int64
	SenderID    *uuid.UUID
	RecipientID *uuid.UUID
	Content     *string
	SentAt      *time.Time
	IsRead      *bool
}

type MessageQueryRepository interface {
	Find(ctx context.Context, messageID int64) (*entity.Message, error)
	Query(ctx context.Context, q *MessageQuery) ([]*entity.Message, error)
}

type MessageCommandRepository interface {
	Create(ctx context.Context, message *entity.Message) (error)
	Update(ctx context.Context, message *entity.Message) error
	Delete(ctx context.Context, messageID int64) error
}

type MessageRepository interface {
	MessageQueryRepository
	MessageCommandRepository
}
