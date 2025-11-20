package domain

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Picture struct {
	ID           int32        `db:"id"`
	UserID       uuid.UUID    `db:"user_id"`
	URL          string       `db:"url"`
	IsProfilePic sql.NullBool `db:"is_profile_pic"`
	CreatedAt    time.Time    `db:"created_at"`
}

type PictureQueryRepository interface {
	Find(ctx context.Context, pictureID int32) (*Picture, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Picture, error)
}

type PictureCommandRepository interface {
	Save(ctx context.Context, picture *Picture) error
	Delete(ctx context.Context, pictureID int32) error
}

type PictureRepository interface {
	PictureQueryRepository
	PictureCommandRepository
}
