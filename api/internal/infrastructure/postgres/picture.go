package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type pictureRepository struct {
	db DBTX
}

func NewPictureRepository(db DBTX) domain.PictureRepository {
	return &pictureRepository{db: db}
}

func (r *pictureRepository) Save(ctx context.Context, picture *domain.Picture) error {
	query := `
		INSERT INTO pictures (id, user_id, url, is_profile_pic, created_at)
		VALUES (:id, :user_id, :url, :is_profile_pic, :created_at)
		ON CONFLICT (id) DO UPDATE SET
			url = :url,
			is_profile_pic = :is_profile_pic
	`
	_, err := r.db.NamedExecContext(ctx, query, picture)
	return err
}

func (r *pictureRepository) Find(ctx context.Context, pictureID int32) (*domain.Picture, error) {
	var picture domain.Picture
	query := "SELECT * FROM pictures WHERE id = $1"
	err := r.db.GetContext(ctx, &picture, query, pictureID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &picture, nil
}

func (r *pictureRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Picture, error) {
	var pictures []domain.Picture
	query := "SELECT * FROM pictures WHERE user_id = $1"
	err := r.db.SelectContext(ctx, &pictures, query, userID)
	if err != nil {
		return nil, err
	}
	return pictures, nil
}

func (r *pictureRepository) Delete(ctx context.Context, pictureID int32) error {
	query := "DELETE FROM pictures WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, pictureID)
	return err
}
