package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type blockRepository struct {
	db DBTX
}

func NewBlockRepository(db DBTX) domain.BlockRepository {
	return &blockRepository{db: db}
}

func (r *blockRepository) Save(ctx context.Context, block *domain.Block) error {
	query := `
		INSERT INTO blocks (blocker_id, blocked_id)
		VALUES (:blocker_id, :blocked_id)
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, block)
	return err
}

func (r *blockRepository) Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := "DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2"
	_, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	return err
}

func (r *blockRepository) Find(ctx context.Context, blockerID, blockedID uuid.UUID) (*domain.Block, error) {
	var block domain.Block
	query := "SELECT * FROM blocks WHERE blocker_id = $1 AND blocked_id = $2"
	err := r.db.GetContext(ctx, &block, query, blockerID, blockedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

func (r *blockRepository) GetByBlockerID(ctx context.Context, blockerID uuid.UUID) ([]domain.Block, error) {
	var blocks []domain.Block
	query := "SELECT * FROM blocks WHERE blocker_id = $1"
	err := r.db.SelectContext(ctx, &blocks, query, blockerID)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}
