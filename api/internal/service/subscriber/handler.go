package subscriber

import (
	"context"
	"encoding/json"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
	"log"
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
	readPayload := &client.ReadPayload{
		UserID:      payload.UserID,
		RecipientID: payload.RecipientID,
		Timestamp:   payload.Timestamp,
	}
	payloadBytes, err := json.Marshal(readPayload)
	if err != nil {
		return err
	}
	if err := h.readPub.Publish(ctx, payloadBytes); err != nil {
		return err
	}
	return nil
}

func (h *subscriberHandler) ChatSubscHandler(ctx context.Context, payload *client.MessagePayload) error {
	log.Printf("Received message payload: %+v", payload)
	msg := &entity.Message{
		SenderID:    payload.SenderID,
		RecipientID: payload.RecipientID,
		Content:     payload.Content,
		SentAt:      payload.SentAt,
	}
	if err := h.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		log.Printf("Creating message from %s to %s: %s", msg.SenderID, msg.RecipientID, msg.Content)
		if err := rm.MessageRepo().Create(ctx, msg); err != nil {
			return err
		}
		log.Printf("Created message with ID %d", msg.ID)
		return nil
	}); err != nil {
		return err
	}
	ackPayload := &client.AckPayload{
		UserID:    msg.SenderID,
		MessageID: msg.ID,
		Timestamp: time.Now().UnixMilli(),
	}
	ackBytes, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}
	if err := h.ackPub.Publish(ctx, ackBytes); err != nil {
		return err
	}

	chatPayload := &client.MessagePayload{
		ID:          msg.ID,
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Content:     msg.Content,
		SentAt:      msg.SentAt,
	}
	chatBytes, err := json.Marshal(chatPayload)
	if err != nil {
		return err
	}
	if err := h.chatPub.Publish(ctx, chatBytes); err != nil {
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
		presencePayload := &client.PresencePayload{
			UserID:      payload.UserID,
			Status:      payload.Status,
			RecipientID: recipientID,
		}
		payloadBytes, err := json.Marshal(presencePayload)
		if err != nil {
			// In a loop, maybe just log and continue
			log.Printf("Error marshalling presence payload: %v", err)
			continue
		}
		if err := h.presencePub.Publish(ctx, payloadBytes); err != nil {
			log.Printf("Error publishing presence payload: %v", err)
			// Decide if we should continue or return
		}
	}
	return nil
}
