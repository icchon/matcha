package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type PictureQuery struct {
	ID           *int32
	UserID       *uuid.UUID
	URL          *string
	IsProfilePic *bool
	CreatedAt    *time.Time
}

type CreatePictureParams struct {
	UserID       uuid.UUID
	URL          string
	IsProfilePic sql.NullBool
}

type PictureQueryRepository interface {
	Find(ctx context.Context, pictureID int32) (*entity.Picture, error)
	Query(ctx context.Context, q *PictureQuery) ([]*entity.Picture, error)
}

type PictureCommandRepository interface {
	Create(ctx context.Context, params CreatePictureParams) (*entity.Picture, error)
	Update(ctx context.Context, picture *entity.Picture) error
	Delete(ctx context.Context, pictureID int32) error
}

type PictureRepository interface {
	PictureQueryRepository
	PictureCommandRepository
}
