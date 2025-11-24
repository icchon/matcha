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

type viewRepository struct {
	db DBTX
}

func NewViewRepository(db DBTX) repo.ViewRepository {
	return &viewRepository{db: db}
}

func (r *viewRepository) Create(ctx context.Context, view *entity.View) error {
	query := `
		INSERT INTO views (viewer_id, viewed_id, view_time)
		VALUES (:viewer_id, :viewed_id, :view_time)
	`
	_, err := r.db.NamedExecContext(ctx, query, view)
	return err
}

func (r *viewRepository) Delete(ctx context.Context, viewerID, viewedID uuid.UUID) error {
	query := "DELETE FROM views WHERE viewer_id = $1 AND viewed_id = $2"
	_, err := r.db.ExecContext(ctx, query, viewerID, viewedID)
	return err
}

func (r *viewRepository) Find(ctx context.Context, viewerID, viewedID uuid.UUID) (*entity.View, error) {
	var view entity.View
	query := "SELECT * FROM views WHERE viewer_id = $1 AND viewed_id = $2"
	err := r.db.GetContext(ctx, &view, query, viewerID, viewedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &view, nil
}

func (r *viewRepository) Query(ctx context.Context, q *repo.ViewQuery) ([]*entity.View, error) {
	query := "SELECT * FROM views WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.ViewerID != nil {
		query += fmt.Sprintf(" AND viewer_id = $%d", argCount)
		args = append(args, *q.ViewerID)
		argCount++
	}
	if q.ViewedID != nil {
		query += fmt.Sprintf(" AND viewed_id = $%d", argCount)
		args = append(args, *q.ViewedID)
		argCount++
	}
	if q.ViewTime != nil {
		query += fmt.Sprintf(" AND view_time = $%d", argCount)
		args = append(args, *q.ViewTime)
		argCount++
	}

	var views []*entity.View
	if err := r.db.SelectContext(ctx, &views, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return views, nil
}
