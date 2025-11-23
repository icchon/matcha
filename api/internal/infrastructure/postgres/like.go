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

type likeRepository struct {
	db DBTX
}

func NewLikeRepository(db DBTX) repo.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, like *entity.Like) error {
	query := `
		INSERT INTO likes (liker_id, liked_id)
		VALUES (:liker_id, :liked_id)
		ON CONFLICT (liker_id, liked_id) DO NOTHING
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, like).StructScan(like)
}

func (r *likeRepository) Delete(ctx context.Context, likerID, likedID uuid.UUID) error {
	query := "DELETE FROM likes WHERE liker_id = $1 AND liked_id = $2"
	_, err := r.db.ExecContext(ctx, query, likerID, likedID)
	return err
}

func (r *likeRepository) Find(ctx context.Context, likerID, likedID uuid.UUID) (*entity.Like, error) {
	var like entity.Like
	query := "SELECT * FROM likes WHERE liker_id = $1 AND liked_id = $2"
	err := r.db.GetContext(ctx, &like, query, likerID, likedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) Query(ctx context.Context, q *repo.LikeQuery) ([]*entity.Like, error) {
	query := "SELECT * FROM likes WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.LikerID != nil {
		query += fmt.Sprintf(" AND liker_id = $%d", argCount)
		args = append(args, *q.LikerID)
		argCount++
	}
	if q.LikedID != nil {
		query += fmt.Sprintf(" AND liked_id = $%d", argCount)
		args = append(args, *q.LikedID)
		argCount++
	}
	if q.CreatedAt != nil {
		query += fmt.Sprintf(" AND created_at = $%d", argCount)
		args = append(args, *q.CreatedAt)
		argCount++
	}

	var likes []*entity.Like
	if err := r.db.SelectContext(ctx, &likes, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return likes, nil
}
