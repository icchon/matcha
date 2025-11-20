package domain

import (
	"context"
	"github.com/google/uuid"
)

type Block struct {
	BlockerID uuid.UUID `db:"blocker_id"`
	BlockedID uuid.UUID `db:"blocked_id"`
}

type BlockQueryRepository interface {
	Find(ctx context.Context, blockerID, blockedID uuid.UUID) (*Block, error)
	GetByBlockerID(ctx context.Context, blockerID uuid.UUID) ([]Block, error)
}

type BlockCommandRepository interface {
	Save(ctx context.Context, block *Block) error
	Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error
}

type BlockRepository interface {
	BlockQueryRepository
	BlockCommandRepository
}
