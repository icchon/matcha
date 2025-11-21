package entity

import (
	"github.com/google/uuid"
	"time"
)

type View struct {
	ViewerID uuid.UUID `db:"viewer_id"`
	ViewedID uuid.UUID `db:"viewed_id"`
	ViewTime time.Time `db:"view_time"`
}
