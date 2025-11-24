package profile

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// mockRepositoryManager is a mock for repo.RepositoryManager.
type mockRepositoryManager struct {
	repo.RepositoryManager // Embed interface to avoid implementing all methods
	profileRepo            repo.UserProfileRepository
	viewRepo               repo.ViewRepository
}

func (m *mockRepositoryManager) ProfileRepo() repo.UserProfileRepository {
	return m.profileRepo
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

func TestProfileService_FindProfile(t *testing.T) {
	userID := uuid.New()
	expectedProfile := &entity.UserProfile{
		UserID:   userID,
		Username: sql.NullString{String: "testuser", Valid: true},
	}
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(profileRepo *mock.MockUserProfileRepository)
		userID      uuid.UUID
		expected    *entity.UserProfile
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(expectedProfile, nil)
			},
			userID:      userID,
			expected:    expectedProfile,
			expectedErr: nil,
		},
		{
			name: "Not Found",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, nil)
			},
			userID:      userID,
			expected:    nil,
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "DB Error",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, dbErr)
			},
			userID:      userID,
			expected:    nil,
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(profileRepo)
			}

			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(nil, profileRepo, fileClient, pictureRepo, viewRepo, likeRepo)

			profile, err := service.FindProfile(context.Background(), tc.userID)

			assert.Equal(t, tc.expected, profile)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestProfileService_CreateProfile(t *testing.T) {
	userID := uuid.New()
	profileToCreate := &entity.UserProfile{
		UserID:   userID,
		Username: sql.NullString{String: "newuser", Valid: true},
	}
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(profileRepo *mock.MockUserProfileRepository)
		uowError    error
		profile     *entity.UserProfile
		expected    *entity.UserProfile
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Create(gomock.Any(), profileToCreate).Return(nil)
			},
			profile:     profileToCreate,
			expected:    profileToCreate,
			expectedErr: nil,
		},
		{
			name: "DB Error on Create",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Create(gomock.Any(), profileToCreate).Return(dbErr)
			},
			profile:     profileToCreate,
			expected:    nil,
			expectedErr: apperrors.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(profileRepo)
			}

			mockRM := &mockRepositoryManager{profileRepo: profileRepo}
			mockUOW := &mockUow{rm: mockRM, err: tc.uowError}

			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(mockUOW, profileRepo, fileClient, pictureRepo, viewRepo, likeRepo)

			profile, err := service.CreateProfile(context.Background(), tc.profile)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, profile)
		})
	}
}

func TestProfileService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	originalProfile := &entity.UserProfile{
		UserID:    userID,
		Username:  sql.NullString{String: "original", Valid: true},
		FirstName: sql.NullString{String: "Original", Valid: true},
	}
	updatePayload := &entity.UserProfile{
		Username: sql.NullString{String: "updated", Valid: true},
	}
	// Note: a new object is created inside UpdateProfile, so we can't compare pointers
	expectedProfile := &entity.UserProfile{
		UserID:    userID,
		Username:  sql.NullString{String: "updated", Valid: true},
		FirstName: sql.NullString{String: "Original", Valid: true},
	}
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(profileRepo *mock.MockUserProfileRepository)
		uowError    error
		expected    *entity.UserProfile
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				// A new object is created in the method, so we need to return a copy
				profileCopy := *originalProfile
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(&profileCopy, nil)
				profileRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected:    expectedProfile,
			expectedErr: nil,
		},
		{
			name: "Find returns Not Found",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, nil)
			},
			expected:    nil,
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "Find returns DB Error",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, dbErr)
			},
			expected:    nil,
			expectedErr: apperrors.ErrInternalServer,
		},
		{
			name: "Update returns DB Error",
			setupMocks: func(profileRepo *mock.MockUserProfileRepository) {
				profileCopy := *originalProfile
				profileRepo.EXPECT().Find(gomock.Any(), userID).Return(&profileCopy, nil)
				profileRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			expected:    nil,
			expectedErr: apperrors.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(profileRepo)
			}

			mockRM := &mockRepositoryManager{profileRepo: profileRepo}
			mockUOW := &mockUow{rm: mockRM, err: tc.uowError}

			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(mockUOW, profileRepo, fileClient, pictureRepo, viewRepo, likeRepo)

			profile, err := service.UpdateProfile(context.Background(), userID, updatePayload)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, profile)
		})
	}
}

func TestProfileService_VeiwProfile(t *testing.T) {
	viewerID := uuid.New()
	viewedID := uuid.New()
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(viewRepo *mock.MockViewRepository)
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(viewRepo *mock.MockViewRepository) {
				viewRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "DB Error",
			setupMocks: func(viewRepo *mock.MockViewRepository) {
				viewRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			viewRepo := mock.NewMockViewRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(viewRepo)
			}

			mockRM := &mockRepositoryManager{viewRepo: viewRepo}
			mockUOW := &mockUow{rm: mockRM}

			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewQueryRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(mockUOW, profileRepo, fileClient, pictureRepo, viewQueryRepo, likeRepo)

			err := service.VeiwProfile(context.Background(), viewerID, viewedID)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestProfileService_FindWhoViewedMeList(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedViews := []*entity.View{
		{ViewerID: uuid.New(), ViewedID: userID},
	}

	testCases := []struct {
		name        string
		setupMocks  func(viewRepo *mock.MockViewQueryRepository)
		expected    []*entity.View
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(viewRepo *mock.MockViewQueryRepository) {
				viewRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(expectedViews, nil)
			},
			expected:    expectedViews,
			expectedErr: nil,
		},
		{
			name: "DB Error",
			setupMocks: func(viewRepo *mock.MockViewQueryRepository) {
				viewRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			},
			expected:    nil,
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(viewRepo)
			}

			// Other mocks
			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(nil, profileRepo, fileClient, pictureRepo, viewRepo, likeRepo)

			views, err := service.FindWhoViewedMeList(context.Background(), userID)

			assert.Equal(t, tc.expected, views)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestProfileService_FindWhoLikedMeList(t *testing.T) {
	userID := uuid.New()
	dbErr := errors.New("db error")
	expectedLikes := []*entity.Like{
		{LikerID: uuid.New(), LikedID: userID},
	}

	testCases := []struct {
		name        string
		setupMocks  func(likeRepo *mock.MockLikeQueryRepository)
		expected    []*entity.Like
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(likeRepo *mock.MockLikeQueryRepository) {
				likeRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(expectedLikes, nil)
			},
			expected:    expectedLikes,
			expectedErr: nil,
		},
		{
			name: "DB Error",
			setupMocks: func(likeRepo *mock.MockLikeQueryRepository) {
				likeRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, dbErr)
			},
			expected:    nil,
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(likeRepo)
			}

			// Other mocks
			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			pictureRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(nil, profileRepo, fileClient, pictureRepo, viewRepo, likeRepo)

			likes, err := service.FindWhoLikedMeList(context.Background(), userID)

			assert.Equal(t, tc.expected, likes)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
