package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type Chat struct {
	OtherUser   entity.UserProfile `json:"other_user"`
	LastMessage *entity.Message    `json:"last_message"` // Can be nil if no messages yet
}

type ChatService interface {
	GetChatsForUser(ctx context.Context, userID uuid.UUID) ([]*Chat, error)
	GetChatMessages(ctx context.Context, params *GetChatMessagesParams) ([]*entity.Message, error)
}

type GetChatMessagesParams struct {
	UserID1 uuid.UUID
	UserID2 uuid.UUID
	Limit   int
	Offset  int
}
