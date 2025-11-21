package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type BlockQuery struct {
	BlockerID *uuid.UUID
	BlockedID *uuid.UUID
}

type BlockQueryRepository interface {
	Find(ctx context.Context, blockerID, blockedID uuid.UUID) (*entity.Block, error)
	Query(ctx context.Context, q *BlockQuery) ([]*entity.Block, error)
}

type BlockCommandRepository interface {
	Create(ctx context.Context, block *entity.Block) error
	Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error
}

type BlockRepository interface {
	BlockQueryRepository
	BlockCommandRepository
}
