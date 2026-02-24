package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type reportRepository struct {
	db DBTX
}

func NewReportRepository(db DBTX) repo.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(ctx context.Context, report *entity.Report) error {
	query := `
		INSERT INTO reports (reporter_id, reported_id, reason)
		VALUES (:reporter_id, :reported_id, :reason)
		RETURNING id, created_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowxContext(ctx, report).StructScan(report)
}

func (r *reportRepository) Find(ctx context.Context, reportID int64) (*entity.Report, error) {
	var report entity.Report
	query := "SELECT * FROM reports WHERE id = $1"
	err := r.db.GetContext(ctx, &report, query, reportID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) Query(ctx context.Context, q *repo.ReportQuery) ([]*entity.Report, error) {
	query := "SELECT * FROM reports WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.ID != nil {
		query += fmt.Sprintf(" AND id = $%d", argCount)
		args = append(args, *q.ID)
		argCount++
	}
	if q.ReporterID != nil {
		query += fmt.Sprintf(" AND reporter_id = $%d", argCount)
		args = append(args, *q.ReporterID)
		argCount++
	}
	if q.ReportedID != nil {
		query += fmt.Sprintf(" AND reported_id = $%d", argCount)
		args = append(args, *q.ReportedID)
		argCount++
	}

	var reports []*entity.Report
	if err := r.db.SelectContext(ctx, &reports, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return reports, nil
}
