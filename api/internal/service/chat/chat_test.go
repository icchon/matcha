package chat

//
// import (
// "context"
// "testing"
//
// "go.uber.org/mock/gomock"
// "github.com/google/uuid"
// "github.com/icchon/matcha/api/internal/domain/entity"
// "github.com/icchon/matcha/api/internal/domain/service"
// "github.com/icchon/matcha/api/internal/mock"
// "github.com/stretchr/testify/assert"
// )

// func TestChatService_GetChatsForUser(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	connRepo := mock.NewMockConnectionQueryRepository(ctrl)
// 	msgRepo := mock.NewMockMessageQueryRepository(ctrl)
// 	profileSvc := mock.NewMockProfileService(ctrl)

// 	chatSvc := NewChatService(connRepo, msgRepo, profileSvc)

// 	userID := uuid.New()
// 	otherUserID := uuid.New()

// 	connections := []*entity.Connection{
// 		{User1ID: userID, User2ID: otherUserID},
// 	}
// 	latestMsg := &entity.Message{Content: "hello"}
// 	otherUserProfile := &entity.UserProfile{UserID: otherUserID}

// 	connRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(connections, nil)
// 	msgRepo.EXPECT().GetLatest(gomock.Any(), userID, otherUserID).Return(latestMsg, nil)
// 	profileSvc.EXPECT().FindProfile(gomock.Any(), otherUserID).Return(otherUserProfile, nil)

// 	chats, err := chatSvc.GetChatsForUser(context.Background(), userID)

// 	assert.NoError(t, err)
// 	assert.Len(t, chats, 1)
// 	assert.Equal(t, otherUserID, chats[0].OtherUser.UserID)
// 	assert.Equal(t, "hello", chats[0].LastMessage.Content)
// }

// func TestChatService_GetChatMessages(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	connRepo := mock.NewMockConnectionQueryRepository(ctrl)
// 	msgRepo := mock.NewMockMessageQueryRepository(ctrl)
// 	profileSvc := mock.NewMockProfileService(ctrl)

// 	chatSvc := NewChatService(connRepo, msgRepo, profileSvc)

// 	userID1 := uuid.New()
// 	userID2 := uuid.New()
// 	params := &service.GetChatMessagesParams{
// 		UserID1: userID1,
// 		UserID2: userID2,
// 		Limit:   10,
// 		Offset:  0,
// 	}

// 	messages := []*entity.Message{
// 		{Content: "hello"},
// 		{Content: "world"},
// 	}

// 	msgRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(messages, nil)

// 	result, err := chatSvc.GetChatMessages(context.Background(), params)

// 	assert.NoError(t, err)
// 	assert.Len(t, result, 2)
// }
