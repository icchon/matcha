package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// mockRepositoryManager is a mock for repo.RepositoryManager.
type mockRepositoryManager struct {
	repo.RepositoryManager // Embed interface to avoid implementing all methods
	userRepo               repo.UserRepository
	blockRepo              repo.BlockRepository
	connectionRepo         repo.ConnectionRepository
	likeRepo               repo.LikeRepository
	viewRepo               repo.ViewRepository
}

func (m *mockRepositoryManager) UserRepo() repo.UserRepository {
	return m.userRepo
}
func (m *mockRepositoryManager) BlockRepo() repo.BlockRepository {
	return m.blockRepo
}
func (m *mockRepositoryManager) ConnectionRepo() repo.ConnectionRepository {
	return m.connectionRepo
}
func (m *mockRepositoryManager) LikeRepo() repo.LikeRepository {
	return m.likeRepo
}
func (m *mockRepositoryManager) ViewRepo() repo.ViewRepository {
	return m.viewRepo
}

// mockUow is a mock for repo.UnitOfWork for testing services.
type mockUow struct {
	rm  repo.RepositoryManager
	err error
}

func (u *mockUow) Do(ctx context.Context, fn func(rm repo.RepositoryManager) error) error {
	if u.err != nil {
		return u.err
	}
	return fn(u.rm)
}

func TestUserService_DeleteUser(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(userRepo *mock.MockUserRepository)
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(userRepo *mock.MockUserRepository) {
				userRepo.EXPECT().Delete(gomock.Any(), userID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "DB Error",
			setupMocks: func(userRepo *mock.MockUserRepository) {
				userRepo.EXPECT().Delete(gomock.Any(), userID).Return(dbErr)
			},
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo)
			}

			mockRM := &mockRepositoryManager{userRepo: userRepo}
			mockUOW := &mockUow{rm: mockRM}

			service := &userService{uow: mockUOW}
			err := service.DeleteUser(context.Background(), userID)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUserService_BlockUser(t *testing.T) {
	blockerID := uuid.New()
	blockedID := uuid.New()
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(connRepo *mock.MockConnectionRepository, likeRepo *mock.MockLikeRepository, viewRepo *mock.MockViewRepository, blockRepo *mock.MockBlockRepository)
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(connRepo *mock.MockConnectionRepository, likeRepo *mock.MockLikeRepository, viewRepo *mock.MockViewRepository, blockRepo *mock.MockBlockRepository) {
				connRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				likeRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				likeRepo.EXPECT().Delete(gomock.Any(), blockedID, blockerID).Return(nil)
				viewRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				viewRepo.EXPECT().Delete(gomock.Any(), blockedID, blockerID).Return(nil)
				blockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "ConnectionRepo Delete fails",
			setupMocks: func(connRepo *mock.MockConnectionRepository, likeRepo *mock.MockLikeRepository, viewRepo *mock.MockViewRepository, blockRepo *mock.MockBlockRepository) {
				connRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(dbErr)
			},
			expectedErr: dbErr,
		},
		{
			name: "BlockRepo Create fails",
			setupMocks: func(connRepo *mock.MockConnectionRepository, likeRepo *mock.MockLikeRepository, viewRepo *mock.MockViewRepository, blockRepo *mock.MockBlockRepository) {
				connRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				likeRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				likeRepo.EXPECT().Delete(gomock.Any(), blockedID, blockerID).Return(nil)
				viewRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
				viewRepo.EXPECT().Delete(gomock.Any(), blockedID, blockerID).Return(nil)
				blockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			connRepo := mock.NewMockConnectionRepository(ctrl)
			likeRepo := mock.NewMockLikeRepository(ctrl)
			viewRepo := mock.NewMockViewRepository(ctrl)
			blockRepo := mock.NewMockBlockRepository(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(connRepo, likeRepo, viewRepo, blockRepo)
			}

			mockRM := &mockRepositoryManager{
				connectionRepo: connRepo,
				likeRepo:       likeRepo,
				viewRepo:       viewRepo,
				blockRepo:      blockRepo,
			}
			mockUOW := &mockUow{rm: mockRM}

			service := &userService{uow: mockUOW}
			err := service.BlockUser(context.Background(), blockerID, blockedID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUserService_UnblockUser(t *testing.T) {
	blockerID := uuid.New()
	blockedID := uuid.New()
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(blockRepo *mock.MockBlockRepository)
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(blockRepo *mock.MockBlockRepository) {
				blockRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "DB Error",
			setupMocks: func(blockRepo *mock.MockBlockRepository) {
				blockRepo.EXPECT().Delete(gomock.Any(), blockerID, blockedID).Return(dbErr)
			},
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blockRepo := mock.NewMockBlockRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(blockRepo)
			}

			mockRM := &mockRepositoryManager{blockRepo: blockRepo}
			mockUOW := &mockUow{rm: mockRM}

			service := &userService{uow: mockUOW}
			err := service.UnblockUser(context.Background(), blockerID, blockedID)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUserService_LikeUser(t *testing.T) {
	likerID := uuid.New()
	likedID := uuid.New()
	dbErr := errors.New("db error")

	user1ID, user2ID := likerID, likedID
	if likerID.String() > likedID.String() {
		user1ID, user2ID = likedID, likerID
	}

	testCases := []struct {
		name               string
		setupMocks         func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository, notifSvcMock *mock.MockNotificationService)
		isMatch            bool
		expectedConnection *entity.Connection
		expectedErr        error
	}{
		{
			name: "Success - First Like (No Match)",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository, notifSvcMock *mock.MockNotificationService) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				notifSvcMock.EXPECT().CreateAndSendNotification(gomock.Any(), likerID, likedID, entity.NotifLike).Return(nil, nil)
			},
			isMatch:            false,
			expectedConnection: nil,
			expectedErr:        nil,
		},
		{
			name: "Success - Second Like (Match)",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository, notifSvcMock *mock.MockNotificationService) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(&entity.Like{}, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				connRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				notifSvcMock.EXPECT().CreateAndSendNotification(gomock.Any(), likerID, likedID, entity.NotifLike).Return(nil, nil)
				notifSvcMock.EXPECT().CreateAndSendNotification(gomock.Any(), likerID, likedID, entity.NotifMatch).Return(nil, nil)
				notifSvcMock.EXPECT().CreateAndSendNotification(gomock.Any(), likedID, likerID, entity.NotifMatch).Return(nil, nil)
			},
			isMatch:            true,
			expectedConnection: &entity.Connection{User1ID: user1ID, User2ID: user2ID},
			expectedErr:        nil,
		},
		{
			name: "Find Fails",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository, notifSvcMock *mock.MockNotificationService) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, dbErr)
			},
			isMatch:            false,
			expectedConnection: nil,
			expectedErr:        dbErr,
		},
		{
			name: "Like Create Fails",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository, notifSvcMock *mock.MockNotificationService) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			isMatch:            false,
			expectedConnection: nil,
			expectedErr:        dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mock.NewMockLikeRepository(ctrl)
			likeQueryRepo := mock.NewMockLikeQueryRepository(ctrl)
			connRepo := mock.NewMockConnectionRepository(ctrl)
			notifSvc := mock.NewMockNotificationService(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(likeRepo, likeQueryRepo, connRepo, notifSvc)
			}

			mockRM := &mockRepositoryManager{likeRepo: likeRepo, connectionRepo: connRepo}
			mockUOW := &mockUow{rm: mockRM}

			service := &userService{uow: mockUOW, likeRepo: likeQueryRepo, notifSvc: notifSvc}

			conn, err := service.LikeUser(context.Background(), likerID, likedID)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedConnection != nil && conn != nil {
				// Don't compare created_at
				assert.Equal(t, tc.expectedConnection.User1ID, conn.User1ID)
				assert.Equal(t, tc.expectedConnection.User2ID, conn.User2ID)
			} else {
				assert.Equal(t, tc.expectedConnection, conn)
			}
		})
	}
}

func TestUserService_FindBlockList(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedBlocks := []*entity.Block{{BlockerID: userID, BlockedID: uuid.New()}}

	testCases := []struct {
		name           string
		setupMocks     func(blockRepo *mock.MockBlockRepository)
		uowError       error
		expectedBlocks []*entity.Block
		expectedErr    error
	}{
		{
			name: "Success",
			setupMocks: func(blockRepo *mock.MockBlockRepository) {
				blockRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(expectedBlocks, nil)
			},
			expectedBlocks: expectedBlocks,
			expectedErr:    nil,
		},
		{
			name: "Query returns error",
			setupMocks: func(blockRepo *mock.MockBlockRepository) {
				blockRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			},
			expectedBlocks: nil,
			expectedErr:    dbErr,
		},
		{
			name:           "UOW returns error",
			uowError:       dbErr,
			expectedBlocks: nil,
			expectedErr:    dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blockRepo := mock.NewMockBlockRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(blockRepo)
			}
			mockRM := &mockRepositoryManager{blockRepo: blockRepo}
			mockUOW := &mockUow{rm: mockRM, err: tc.uowError}

			service := &userService{uow: mockUOW}
			blocks, err := service.FindBlockList(context.Background(), userID)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedBlocks, blocks)
		})
	}
}

func TestUserService_UnlikeUser(t *testing.T) {
	likerID := uuid.New()
	likedID := uuid.New()
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(likeRepo *mock.MockLikeRepository, connRepo *mock.MockConnectionRepository, notifServiceMock *mock.MockNotificationService)
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(likeRepo *mock.MockLikeRepository, connRepo *mock.MockConnectionRepository, notifServiceMock *mock.MockNotificationService) {
				likeRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(nil)
				connRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(nil)
				notifServiceMock.EXPECT().CreateAndSendNotification(gomock.Any(), likerID, likedID, entity.NotifUnlike).Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			name: "LikeRepo Delete fails",
			setupMocks: func(likeRepo *mock.MockLikeRepository, connRepo *mock.MockConnectionRepository, notifServiceMock *mock.MockNotificationService) {
				likeRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(dbErr)
			},
			expectedErr: dbErr,
		},
		{
			name: "ConnRepo Delete fails",
			setupMocks: func(likeRepo *mock.MockLikeRepository, connRepo *mock.MockConnectionRepository, notifServiceMock *mock.MockNotificationService) {
				likeRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(nil)
				connRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(dbErr)
			},
			expectedErr: dbErr,
		},
		{
			name: "Notification Service fails",
			setupMocks: func(likeRepo *mock.MockLikeRepository, connRepo *mock.MockConnectionRepository, notifServiceMock *mock.MockNotificationService) {
				likeRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(nil)
				connRepo.EXPECT().Delete(gomock.Any(), likerID, likedID).Return(nil)
				notifServiceMock.EXPECT().CreateAndSendNotification(gomock.Any(), likerID, likedID, entity.NotifUnlike).Return(nil, dbErr)
			},
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mock.NewMockLikeRepository(ctrl)
			connRepo := mock.NewMockConnectionRepository(ctrl)
			notifService := mock.NewMockNotificationService(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(likeRepo, connRepo, notifService)
			}

			mockRM := &mockRepositoryManager{likeRepo: likeRepo, connectionRepo: connRepo}
			mockUOW := &mockUow{rm: mockRM}
			service := NewUserService(mockUOW, nil, nil, nil, notifService, nil, nil, nil) // Pass notifService
			err := service.UnlikeUser(context.Background(), likerID, likedID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUserService_FindMyLikedList(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedLikes := []*entity.Like{{LikerID: userID}}

	testCases := []struct {
		name          string
		setupMocks    func(likeQueryRepo *mock.MockLikeQueryRepository)
		expectedLikes []*entity.Like
		expectedErr   error
	}{
		{
			name: "Success",
			setupMocks: func(likeQueryRepo *mock.MockLikeQueryRepository) {
				likeQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(expectedLikes, nil)
			},
			expectedLikes: expectedLikes,
			expectedErr:   nil,
		},
		{
			name: "Query fails",
			setupMocks: func(likeQueryRepo *mock.MockLikeQueryRepository) {
				likeQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			},
			expectedLikes: nil,
			expectedErr:   dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeQueryRepo := mock.NewMockLikeQueryRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(likeQueryRepo)
			}
			service := &userService{likeRepo: likeQueryRepo}
			likes, err := service.FindMyLikedList(context.Background(), userID)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedLikes, likes)
		})
	}
}

func TestUserService_FindMyViewedList(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedViews := []*entity.View{{ViewerID: userID}}

	testCases := []struct {
		name          string
		setupMocks    func(viewQueryRepo *mock.MockViewQueryRepository)
		expectedViews []*entity.View
		expectedErr   error
	}{
		{
			name: "Success",
			setupMocks: func(viewQueryRepo *mock.MockViewQueryRepository) {
				viewQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(expectedViews, nil)
			},
			expectedViews: expectedViews,
			expectedErr:   nil,
		},
		{
			name: "Query fails",
			setupMocks: func(viewQueryRepo *mock.MockViewQueryRepository) {
				viewQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			},
			expectedViews: nil,
			expectedErr:   dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			viewQueryRepo := mock.NewMockViewQueryRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(viewQueryRepo)
			}
			service := &userService{viewRepo: viewQueryRepo}
			views, err := service.FindMyViewedList(context.Background(), userID)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedViews, views)
		})
	}
}

func TestUserService_FindConnections(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedConns := []*entity.Connection{{User1ID: userID}}

	testCases := []struct {
		name          string
		setupMocks    func(connQueryRepo *mock.MockConnectionQueryRepository)
		expectedConns []*entity.Connection
		expectedErr   error
	}{
		{
			name: "Success",
			setupMocks: func(connQueryRepo *mock.MockConnectionQueryRepository) {
				connQueryRepo.EXPECT().Query(gomock.Any(), &repo.ConnectionQuery{User1ID: &userID}).Return(expectedConns, nil)
			},
			expectedConns: expectedConns,
			expectedErr:   nil,
		},
		{
			name: "Query fails",
			setupMocks: func(connQueryRepo *mock.MockConnectionQueryRepository) {
				connQueryRepo.EXPECT().Query(gomock.Any(), &repo.ConnectionQuery{User1ID: &userID}).Return(nil, dbErr)
			},
			expectedConns: nil,
			expectedErr:   dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			connQueryRepo := mock.NewMockConnectionQueryRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(connQueryRepo)
			}
			service := &userService{connectionRepo: connQueryRepo}
			conns, err := service.FindConnections(context.Background(), userID)
			assert.Equal(t, tc.expectedErr, err)
			assert.ElementsMatch(t, tc.expectedConns, conns)
		})
	}
}
