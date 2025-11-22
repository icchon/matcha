package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
)

// mockRepositoryManager is a mock for uow.RepositoryManager.
type mockRepositoryManager struct {
	uow.RepositoryManager // Embed interface to avoid implementing all methods
	userRepo       repo.UserRepository
	blockRepo      repo.BlockRepository
	connectionRepo repo.ConnectionRepository
	likeRepo       repo.LikeRepository
	viewRepo       repo.ViewRepository
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

// mockUow is a mock for uow.UnitOfWork for testing services.
type mockUow struct {
	rm  uow.RepositoryManager
	err error
}

func (u *mockUow) Do(ctx context.Context, fn func(rm uow.RepositoryManager) error) error {
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

	testCases := []struct {
		name              string
		setupMocks        func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository)
		isMatch           bool
		expectedConnection *entity.Connection
		expectedErr       error
	}{
		{
			name: "Success - First Like (No Match)",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			isMatch:           false,
			expectedConnection: nil,
			expectedErr:       nil,
		},
		{
			name: "Success - Second Like (Match)",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(&entity.Like{}, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				connRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			isMatch:           true,
			expectedConnection: &entity.Connection{User1ID: likerID, User2ID: likedID},
			expectedErr:       nil,
		},
		{
			name: "Find Fails",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, dbErr)
			},
			isMatch:           false,
			expectedConnection: nil,
			expectedErr:       dbErr,
		},
		{
			name: "Like Create Fails",
			setupMocks: func(likeRepoMock *mock.MockLikeRepository, likeQueryRepoMock *mock.MockLikeQueryRepository, connRepoMock *mock.MockConnectionRepository) {
				likeQueryRepoMock.EXPECT().Find(gomock.Any(), likedID, likerID).Return(nil, nil)
				likeRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			isMatch:           false,
			expectedConnection: nil,
			expectedErr:       dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mock.NewMockLikeRepository(ctrl)
			likeQueryRepo := mock.NewMockLikeQueryRepository(ctrl)
			connRepo := mock.NewMockConnectionRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(likeRepo, likeQueryRepo, connRepo)
			}
			
			mockRM := &mockRepositoryManager{likeRepo: likeRepo, connectionRepo: connRepo}
			mockUOW := &mockUow{rm: mockRM}
			
			service := &userService{uow: mockUOW, likeRepo: likeQueryRepo}

			conn, err := service.LikeUser(context.Background(), likerID, likedID)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedConnection, conn)
		})
	}
}