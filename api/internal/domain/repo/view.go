package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type ViewQuery struct {
	ViewerID *uuid.UUID
	ViewedID *uuid.UUID
	ViewTime *time.Time
}

type ViewQueryRepository interface {
	Find(ctx context.Context, viewerID, viewedID uuid.UUID) (*entity.View, error)
	Query(ctx context.Context, q *ViewQuery) ([]*entity.View, error)
}

type ViewCommandRepository interface {
	Save(ctx context.Context, view *entity.View) error
	Delete(ctx context.Context, viewerID, viewedID uuid.UUID) error
}

type ViewRepository interface {
	ViewQueryRepository
	ViewCommandRepository
}
