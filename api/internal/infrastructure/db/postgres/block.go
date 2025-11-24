package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type blockRepository struct {
	db DBTX
}

func NewBlockRepository(db DBTX) repo.BlockRepository {
	return &blockRepository{db: db}
}

func (r *blockRepository) Create(ctx context.Context, block *entity.Block) error {
	query := `
		INSERT INTO blocks (blocker_id, blocked_id)
		VALUES (:blocker_id, :blocked_id)
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, block).StructScan(block)
}

func (r *blockRepository) Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := "DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2"
	_, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	return err
}

func (r *blockRepository) Find(ctx context.Context, blockerID, blockedID uuid.UUID) (*entity.Block, error) {
	var block entity.Block
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

func (r *blockRepository) Query(ctx context.Context, q *repo.BlockQuery) ([]*entity.Block, error) {
	query := "SELECT * FROM blocks WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.BlockerID != nil {
		query += fmt.Sprintf(" AND blocker_id = $%d", argCount)
		args = append(args, *q.BlockerID)
		argCount++
	}
	if q.BlockedID != nil {
		query += fmt.Sprintf(" AND blocked_id = $%d", argCount)
		args = append(args, *q.BlockedID)
		argCount++
	}

	var blocks []*entity.Block
	if err := r.db.SelectContext(ctx, &blocks, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return blocks, nil
}
