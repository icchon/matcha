package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type likeRepository struct {
	db DBTX
}

func NewLikeRepository(db DBTX) domain.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Save(ctx context.Context, like *domain.Like) error {
	query := `
		INSERT INTO likes (liker_id, liked_id, created_at)
		VALUES (:liker_id, :liked_id, :created_at)
		ON CONFLICT (liker_id, liked_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, like)
	return err
}

func (r *likeRepository) Delete(ctx context.Context, likerID, likedID uuid.UUID) error {
	query := "DELETE FROM likes WHERE liker_id = $1 AND liked_id = $2"
	_, err := r.db.ExecContext(ctx, query, likerID, likedID)
	return err
}

func (r *likeRepository) Find(ctx context.Context, likerID, likedID uuid.UUID) (*domain.Like, error) {
	var like domain.Like
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

func (r *likeRepository) GetByLikerID(ctx context.Context, likerID uuid.UUID) ([]domain.Like, error) {
	var likes []domain.Like
	query := "SELECT * FROM likes WHERE liker_id = $1"
	err := r.db.SelectContext(ctx, &likes, query, likerID)
	if err != nil {
		return nil, err
	}
	return likes, nil
}
