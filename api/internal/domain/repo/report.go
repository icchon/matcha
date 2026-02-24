package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type ReportQuery struct {
	ID         *int64
	ReporterID *uuid.UUID
	ReportedID *uuid.UUID
}

type ReportQueryRepository interface {
	Find(ctx context.Context, reportID int64) (*entity.Report, error)
	Query(ctx context.Context, q *ReportQuery) ([]*entity.Report, error)
}

type ReportCommandRepository interface {
	Create(ctx context.Context, report *entity.Report) error
}

type ReportRepository interface {
	ReportQueryRepository
	ReportCommandRepository
}
