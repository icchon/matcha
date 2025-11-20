package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type viewRepository struct {
	db DBTX
}

func NewViewRepository(db DBTX) domain.ViewRepository {
	return &viewRepository{db: db}
}

func (r *viewRepository) Save(ctx context.Context, view *domain.View) error {
	query := `
		INSERT INTO views (viewer_id, viewed_id, view_time)
		VALUES (:viewer_id, :viewed_id, :view_time)
		ON CONFLICT (viewer_id, viewed_id, view_time) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, view)
	return err
}

func (r *viewRepository) Delete(ctx context.Context, viewerID, viewedID uuid.UUID) error {
	query := "DELETE FROM views WHERE viewer_id = $1 AND viewed_id = $2"
	_, err := r.db.ExecContext(ctx, query, viewerID, viewedID)
	return err
}

func (r *viewRepository) Find(ctx context.Context, viewerID, viewedID uuid.UUID) (*domain.View, error) {
	var view domain.View
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

func (r *viewRepository) GetByViewerID(ctx context.Context, viewerID uuid.UUID) ([]domain.View, error) {
	var views []domain.View
	query := "SELECT * FROM views WHERE viewer_id = $1"
	err := r.db.SelectContext(ctx, &views, query, viewerID)
	if err != nil {
		return nil, err
	}
	return views, nil
}
