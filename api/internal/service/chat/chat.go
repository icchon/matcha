package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity" // Added import
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type chatService struct {
	connRepo    repo.ConnectionQueryRepository
	messageRepo repo.MessageQueryRepository
	profileSvc  service.ProfileService
}

var _ service.ChatService = (*chatService)(nil)

func NewChatService(connRepo repo.ConnectionQueryRepository, messageRepo repo.MessageQueryRepository, profileSvc service.ProfileService) *chatService {
	return &chatService{
		connRepo:    connRepo,
		messageRepo: messageRepo,
		profileSvc:  profileSvc,
	}
}

func (s *chatService) GetChatsForUser(ctx context.Context, userID uuid.UUID) ([]*service.Chat, error) {
	q := &repo.ConnectionQuery{
		User1ID: &userID,
	}
	connections, err := s.connRepo.Query(ctx, q)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	chats := make([]*service.Chat, 0, len(connections))

	for _, conn := range connections {
		otherUserID, err := conn.GetOtherUserID(userID)
		if err != nil {
			continue
		}

		latestMessage, err := s.messageRepo.GetLatest(ctx, userID, otherUserID)
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}

		otherUserProfile, err := s.profileSvc.FindProfile(ctx, otherUserID)
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}
		if otherUserProfile == nil {
			continue
		}

		chats = append(chats, &service.Chat{
			OtherUser:   *otherUserProfile,
			LastMessage: latestMessage,
		})
	}

	return chats, nil
}

func (s *chatService) GetChatMessages(ctx context.Context, params *service.GetChatMessagesParams) ([]*entity.Message, error) {
	q := &repo.MessageQuery{
		SenderID:    &params.UserID1,
		RecipientID: &params.UserID2,
		Limit:       &params.Limit,
		Offset:      &params.Offset,
	}

	messages, err := s.messageRepo.Query(ctx, q)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return messages, nil
}
