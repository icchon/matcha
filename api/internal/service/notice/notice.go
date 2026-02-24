package notice

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors" // Added this import
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

func (s *notificationService) MarkNotificationAsRead(ctx context.Context, notificationID int64, recipientID uuid.UUID) error {
	notification, err := s.notificationRepo.Find(ctx, notificationID)
	if err != nil {
		return err // Handle database errors
	}
	if notification == nil {
		return apperrors.ErrNotFound // Notification not found
	}
	if notification.RecipientID != recipientID {
		return apperrors.ErrUnauthorized // User is not the recipient of this notification
	}

	// Only update if it's not already read
	if !notification.IsRead.Bool {
		notification.IsRead = sql.NullBool{Bool: true, Valid: true}
		return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
			return rm.NotificationRepo().Update(ctx, notification)
		})
	}
	return nil // Already read, no action needed
}

func (s *notificationService) MarkAllNotificationsAsRead(ctx context.Context, recipientID uuid.UUID) error {
	isRead := false
	notifications, err := s.notificationRepo.Query(ctx, &repo.NotificationQuery{
		RecipientID: &recipientID,
		IsRead:      &isRead,
	})
	if err != nil {
		return err
	}

	if len(notifications) == 0 {
		return nil // No unread notifications to mark
	}

	return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		for _, notif := range notifications {
			notif.IsRead = sql.NullBool{Bool: true, Valid: true}
			if err := rm.NotificationRepo().Update(ctx, notif); err != nil {
				return err
			}
		}
		return nil
	})
}
