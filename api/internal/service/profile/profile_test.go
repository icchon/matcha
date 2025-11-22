package profile

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
)

// mockRepositoryManager is a mock for uow.RepositoryManager.
type mockRepositoryManager struct {
	uow.RepositoryManager // Embed interface to avoid implementing all methods
	profileRepo           repo.UserProfileRepository
	viewRepo              repo.ViewRepository
}

func (m *mockRepositoryManager) ProfileRepo() repo.UserProfileRepository {
	return m.profileRepo
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
