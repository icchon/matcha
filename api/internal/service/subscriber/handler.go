package subscriber

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
	"time"
)

type subscriberHandler struct {
	uow         repo.UnitOfWork
	messageRepo repo.MessageRepository
	readPub     client.Publisher
	ackPub      client.Publisher
	chatPub     client.Publisher
	presencePub client.Publisher

	userService  service.UserService
	notifService service.NotificationService
}

func NewSubscriberHandler(
	uow repo.UnitOfWork,
	messageRepo repo.MessageRepository,
	readPub client.Publisher,
	ackPub client.Publisher,
	chatPub client.Publisher,
	presencePub client.Publisher,
	userService service.UserService,
	notifService service.NotificationService,
) *subscriberHandler {
	return &subscriberHandler{
		uow:          uow,
		messageRepo:  messageRepo,
		readPub:      readPub,
		ackPub:       ackPub,
		chatPub:      chatPub,
		presencePub:  presencePub,
		userService:  userService,
		notifService: notifService,
	}
}

func (h *subscriberHandler) ReadSubscHandler(ctx context.Context, payload *client.ReadPayload) error {
	unreadMsgs, err := h.messageRepo.Query(ctx, &repo.MessageQuery{
		SenderID:    &payload.RecipientID,
		RecipientID: &payload.UserID,
		IsRead:      func(b bool) *bool { return &b }(false),
	})
	if err != nil {
		return err
	}
	if err := h.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		for _, msg := range unreadMsgs {
			if err := rm.MessageRepo().MarkAsRead(ctx, msg.ID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := h.readPub.Publish(ctx, &client.ReadPayload{
		UserID:      payload.UserID,
		RecipientID: payload.RecipientID,
		Timestamp:   payload.Timestamp,
	}); err != nil {
		return err
	}
	return nil
}

func (h *subscriberHandler) ChatSubscHandler(ctx context.Context, payload *client.MessagePayload) error {
	msg := &entity.Message{
		SenderID:    payload.SenderID,
		RecipientID: payload.RecipientID,
		Content:     payload.Content,
		SentAt:      payload.SentAt,
	}
	if err := h.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.MessageRepo().Create(ctx, msg); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if err := h.ackPub.Publish(ctx, &client.AckPayload{
		UserID:    msg.SenderID,
		MessageID: msg.ID,
		Timestamp: time.Now().UnixMilli(),
	}); err != nil {
		return err
	}
	if err := h.chatPub.Publish(ctx, &client.MessagePayload{
		ID:          msg.ID,
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Content:     msg.Content,
		SentAt:      msg.SentAt,
	}); err != nil {
		return err
	}
	if _, err := h.notifService.CreateAndSendNotofication(ctx, msg.SenderID, msg.RecipientID, entity.NotifMessage); err != nil {
		return err
	}
	return nil
}

func (h *subscriberHandler) PresenceSubscHandler(ctx context.Context, payload *client.PresencePayload) error {
	connections, err := h.userService.FindConnections(ctx, payload.UserID)
	if err != nil {
		return err
	}
	for _, conn := range connections {
		recipientID, err := conn.GetOtherUserID(payload.UserID)
		if err != nil {
			return err
		}
		if err := h.presencePub.Publish(ctx, &client.PresencePayload{
			UserID:      payload.UserID,
			Status:      payload.Status,
			RecipientID: recipientID,
		}); err != nil {
			return err
		}
	}
	return nil
}
