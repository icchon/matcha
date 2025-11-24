package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type pictureRepository struct {
	db DBTX
}

func NewPictureRepository(db DBTX) repo.PictureRepository {
	return &pictureRepository{db: db}
}

func (r *pictureRepository) Create(ctx context.Context, picture *entity.Picture) error {
	query := `
		INSERT INTO pictures (user_id, url, is_profile_pic)
		VALUES (:user_id, :url, :is_profile_pic)
		RETURNING *
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowxContext(ctx, picture).StructScan(picture)
}

func (r *pictureRepository) Update(ctx context.Context, picture *entity.Picture) error {
	query := `
		UPDATE pictures SET
			url = :url,
			is_profile_pic = :is_profile_pic
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, picture)
	return err
}

func (r *pictureRepository) Find(ctx context.Context, pictureID int32) (*entity.Picture, error) {
	var picture entity.Picture
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

func (r *pictureRepository) Query(ctx context.Context, q *repo.PictureQuery) ([]*entity.Picture, error) {
	query := "SELECT * FROM pictures WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.ID != nil {
		query += fmt.Sprintf(" AND id = $%d", argCount)
		args = append(args, *q.ID)
		argCount++
	}
	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.URL != nil {
		query += fmt.Sprintf(" AND url = $%d", argCount)
		args = append(args, *q.URL)
		argCount++
	}
	if q.IsProfilePic != nil {
		query += fmt.Sprintf(" AND is_profile_pic = $%d", argCount)
		args = append(args, *q.IsProfilePic)
		argCount++
	}
	if q.CreatedAt != nil {
		query += fmt.Sprintf(" AND created_at = $%d", argCount)
		args = append(args, *q.CreatedAt)
		argCount++
	}

	var pictures []*entity.Picture
	if err := r.db.SelectContext(ctx, &pictures, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return pictures, nil
}

func (r *pictureRepository) Delete(ctx context.Context, pictureID int32) error {
	query := "DELETE FROM pictures WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, pictureID)
	return err
}
