package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type NotificationService interface {
	GetNotifications(ctx context.Context, recipientID uuid.UUID) ([]*entity.Notification, error)
	CreateAndSendNotofication(ctx context.Context, senderID uuid.UUID, recipiendID uuid.UUID, notifType entity.NotificationType) (*entity.Notification, error)
}
