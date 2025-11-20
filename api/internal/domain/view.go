package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type View struct {
	ViewerID uuid.UUID `db:"viewer_id"`
	ViewedID uuid.UUID `db:"viewed_id"`
	ViewTime time.Time `db:"view_time"`
}

type ViewQueryRepository interface {
	Find(ctx context.Context, viewerID, viewedID uuid.UUID) (*View, error)
	GetByViewerID(ctx context.Context, viewerID uuid.UUID) ([]View, error)
}

type ViewCommandRepository interface {
	Save(ctx context.Context, view *View) error
	Delete(ctx context.Context, viewerID, viewedID uuid.UUID) error
}

type ViewRepository interface {
	ViewQueryRepository
	ViewCommandRepository
}
