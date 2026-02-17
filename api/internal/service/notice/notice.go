package notice

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

var _ service.NotificationService = (*notificationService)(nil)

type notificationService struct {
	uow              repo.UnitOfWork
	notificationRepo repo.NotificationQueryRepository
	notificationPub  client.Publisher
}

func NewNotificationService(uow repo.UnitOfWork, notificationRepo repo.NotificationQueryRepository, notificationPub client.Publisher) service.NotificationService {
	return &notificationService{
		uow:              uow,
		notificationRepo: notificationRepo,
		notificationPub:  notificationPub,
	}
}

func (s *notificationService) GetNotifications(ctx context.Context, recipientID uuid.UUID) ([]*entity.Notification, error) {
	return s.notificationRepo.Query(ctx, &repo.NotificationQuery{
		RecipientID: &recipientID,
	})
}

func (s *notificationService) CreateAndSendNotification(ctx context.Context, senderID uuid.UUID, recipiendID uuid.UUID, notifType entity.NotificationType) (*entity.Notification, error) {
	notification := &entity.Notification{
		RecipientID: recipiendID,
		SenderID:    sql.NullString{String: senderID.String(), Valid: true},
		Type:        notifType,
	}
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.NotificationRepo().Create(ctx, notification); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	payload := client.NotificationPayload{
		ID:          notification.ID,
		RecipientID: notification.RecipientID,
		SenderID:    notification.SenderID.String,
		Type:        string(notification.Type),
		CreatedAt:   notification.CreatedAt,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if err := s.notificationPub.Publish(ctx, payloadBytes); err != nil {
		return nil, err
	}
	return notification, nil
}
