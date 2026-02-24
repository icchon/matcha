package entity

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID         int64     `db:"id"`
	ReporterID uuid.UUID `db:"reporter_id"`
	ReportedID uuid.UUID `db:"reported_id"`
	Reason     string    `db:"reason"`
	CreatedAt  time.Time `db:"created_at"`
}
