package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type notificationRepository struct {
	db DBTX
}

func NewNotificationRepository(db DBTX) repo.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	query := `
		INSERT INTO notifications (recipient_id, sender_id, type)
		VALUES (:recipient_id, :sender_id, :type)
		RETURNING *
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowxContext(ctx, notification).StructScan(notification)
}

func (r *notificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	query := `
		UPDATE notifications SET
			is_read = :is_read
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, notification)
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, notificationID int64) error {
	query := "DELETE FROM notifications WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, notificationID)
	return err
}

func (r *notificationRepository) Find(ctx context.Context, notificationID int64) (*entity.Notification, error) {
	var notification entity.Notification
	query := "SELECT * FROM notifications WHERE id = $1"
	err := r.db.GetContext(ctx, &notification, query, notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Query(ctx context.Context, q *repo.NotificationQuery) ([]*entity.Notification, error) {
	query := "SELECT * FROM notifications WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.ID != nil {
		query += fmt.Sprintf(" AND id = $%d", argCount)
		args = append(args, *q.ID)
		argCount++
	}
	if q.RecipientID != nil {
		query += fmt.Sprintf(" AND recipient_id = $%d", argCount)
		args = append(args, *q.RecipientID)
		argCount++
	}
	if q.SenderID != nil {
		query += fmt.Sprintf(" AND sender_id = $%d", argCount)
		args = append(args, *q.SenderID)
		argCount++
	}
	if q.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, *q.Type)
		argCount++
	}
	if q.IsRead != nil {
		query += fmt.Sprintf(" AND is_read = $%d", argCount)
		args = append(args, *q.IsRead)
		argCount++
	}
	if q.CreatedAt != nil {
		query += fmt.Sprintf(" AND created_at = $%d", argCount)
		args = append(args, *q.CreatedAt)
		argCount++
	}

	var notifications []*entity.Notification
	if err := r.db.SelectContext(ctx, &notifications, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return notifications, nil
}
