package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type notificationRepository struct {
	db DBTX
}

func NewNotificationRepository(db DBTX) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, recipient_id, sender_id, type, is_read, created_at)
		VALUES (:id, :recipient_id, :sender_id, :type, :is_read, :created_at)
		ON CONFLICT (id) DO UPDATE SET
			is_read = :is_read
	`
	_, err := r.db.NamedExecContext(ctx, query, notification)
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, notificationID int64) error {
	query := "DELETE FROM notifications WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, notificationID)
	return err
}

func (r *notificationRepository) Find(ctx context.Context, notificationID int64) (*domain.Notification, error) {
	var notification domain.Notification
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

func (r *notificationRepository) GetByRecipientID(ctx context.Context, recipientID uuid.UUID) ([]domain.Notification, error) {
	var notifications []domain.Notification
	query := "SELECT * FROM notifications WHERE recipient_id = $1"
	err := r.db.SelectContext(ctx, &notifications, query, recipientID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
