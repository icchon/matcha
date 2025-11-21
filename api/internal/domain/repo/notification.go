package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"time"
)

type NotificationQuery struct {
	ID          *int64
	RecipientID *uuid.UUID
	SenderID    *uuid.UUID
	Type        *entity.NotificationType
	IsRead      *bool
	CreatedAt   *time.Time
}

type CreateNotificationParams struct {
	RecipientID uuid.UUID
	SenderID    uuid.UUID
	Type        entity.NotificationType
}

type NotificationQueryRepository interface {
	Find(ctx context.Context, notificationID int64) (*entity.Notification, error)
	Query(ctx context.Context, q *NotificationQuery) ([]*entity.Notification, error)
}

type NotificationCommandRepository interface {
	Create(ctx context.Context, params CreateNotificationParams) (*entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, notificationID int64) error
}

type NotificationRepository interface {
	NotificationQueryRepository
	NotificationCommandRepository
}
