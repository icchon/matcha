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

// Re-defining mock helpers as they are not in a shared file.

// mockPictureRepositoryManager is a mock for uow.RepositoryManager.
type mockPictureRepositoryManager struct {
	uow.RepositoryManager
	pictureRepo repo.PictureRepository
}

func (m *mockPictureRepositoryManager) PictureRepo() repo.PictureRepository {
	return m.pictureRepo
}

// mockPictureUow is a mock for uow.UnitOfWork for testing services.
type mockPictureUow struct {
	rm  uow.RepositoryManager
	err error
}

func (u *mockPictureUow) Do(ctx context.Context, fn func(rm uow.RepositoryManager) error) error {
	if u.err != nil {
		return u.err
	}
	return fn(u.rm)
}

func TestProfileService_DeletePicture(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	pictureID := int32(1)
	dbErr := errors.New("db error")

	testCases := []struct {
		name        string
		setupMocks  func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository)
		uowError    error
		pictureID   int32
		userID      uuid.UUID
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository) {
				picQueryRepo.EXPECT().Find(gomock.Any(), pictureID).Return(&entity.Picture{ID: pictureID, UserID: userID}, nil)
				picRepo.EXPECT().Delete(gomock.Any(), pictureID).Return(nil)
			},
			pictureID:   pictureID,
			userID:      userID,
			expectedErr: nil,
		},
		{
			name: "Find fails",
			setupMocks: func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository) {
				picQueryRepo.EXPECT().Find(gomock.Any(), pictureID).Return(nil, dbErr)
			},
			pictureID:   pictureID,
			userID:      userID,
			expectedErr: apperrors.ErrInternalServer,
		},
		{
			name: "Picture not found",
			setupMocks: func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository) {
				picQueryRepo.EXPECT().Find(gomock.Any(), pictureID).Return(nil, nil)
			},
			pictureID:   pictureID,
			userID:      userID,
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "User is not owner",
			setupMocks: func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository) {
				picQueryRepo.EXPECT().Find(gomock.Any(), pictureID).Return(&entity.Picture{ID: pictureID, UserID: otherUserID}, nil)
			},
			pictureID:   pictureID,
			userID:      userID,
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "Delete fails in UoW",
			setupMocks: func(picQueryRepo *mock.MockPictureQueryRepository, picRepo *mock.MockPictureRepository) {
				picQueryRepo.EXPECT().Find(gomock.Any(), pictureID).Return(&entity.Picture{ID: pictureID, UserID: userID}, nil)
				picRepo.EXPECT().Delete(gomock.Any(), pictureID).Return(dbErr)
			},
			pictureID:   pictureID,
			userID:      userID,
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			picQueryRepo := mock.NewMockPictureQueryRepository(ctrl)
			picRepo := mock.NewMockPictureRepository(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(picQueryRepo, picRepo)
			}

			mockRM := &mockPictureRepositoryManager{pictureRepo: picRepo}
			mockUOW := &mockPictureUow{rm: mockRM, err: tc.uowError}

			// Dummy mocks for other dependencies of profileService
			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)
			fileClient := mock.NewMockFileClient(ctrl)

			service := NewProfileService(mockUOW, profileRepo, fileClient, picQueryRepo, viewRepo, likeRepo)

			err := service.DeletePicture(context.Background(), tc.pictureID, tc.userID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestProfileService_UploadPicture(t *testing.T) {
	userID := uuid.New()
	imageBytes := []byte("test-image")
	imageURL := "http://example.com/image.jpg"
	dbErr := errors.New("db error")
	fileErr := errors.New("file client error")

	testCases := []struct {
		name        string
		setupMocks  func(fileClient *mock.MockFileClient, picRepo *mock.MockPictureRepository)
		image       []byte
		expectedPic *entity.Picture
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func(fileClient *mock.MockFileClient, picRepo *mock.MockPictureRepository) {
				fileClient.EXPECT().SaveImage(imageBytes, gomock.Any()).Return(imageURL, nil)
				picRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *entity.Picture) error {
					// Check if the picture being created has the correct data
					assert.Equal(t, userID, p.UserID)
					assert.Equal(t, imageURL, p.URL)
					assert.False(t, p.IsProfilePic.Bool)
					return nil
				})
			},
			image: imageBytes,
			expectedPic: &entity.Picture{
				UserID: userID,
				URL:    imageURL,
				IsProfilePic: sql.NullBool{Bool: false, Valid: true},
			},
			expectedErr: nil,
		},
		{
			name:        "Empty Image",
			setupMocks:  nil,
			image:       []byte{},
			expectedPic: nil,
			expectedErr: apperrors.ErrInvalidInput,
		},
		{
			name: "FileClient SaveImage fails",
			setupMocks: func(fileClient *mock.MockFileClient, picRepo *mock.MockPictureRepository) {
				fileClient.EXPECT().SaveImage(imageBytes, gomock.Any()).Return("", fileErr)
			},
			image:       imageBytes,
			expectedPic: nil,
			expectedErr: fileErr,
		},
		{
			name: "DB Create fails",
			setupMocks: func(fileClient *mock.MockFileClient, picRepo *mock.MockPictureRepository) {
				fileClient.EXPECT().SaveImage(imageBytes, gomock.Any()).Return(imageURL, nil)
				picRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)
			},
			image:       imageBytes,
			expectedPic: nil,
			expectedErr: dbErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileClient := mock.NewMockFileClient(ctrl)
			picRepo := mock.NewMockPictureRepository(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(fileClient, picRepo)
			}
			
			mockRM := &mockPictureRepositoryManager{pictureRepo: picRepo}
			mockUOW := &mockPictureUow{rm: mockRM}

			// Dummy mocks for other dependencies
			profileRepo := mock.NewMockUserProfileRepository(ctrl)
			picQueryRepo := mock.NewMockPictureQueryRepository(ctrl)
			viewRepo := mock.NewMockViewQueryRepository(ctrl)
			likeRepo := mock.NewMockLikeQueryRepository(ctrl)

			service := NewProfileService(mockUOW, profileRepo, fileClient, picQueryRepo, viewRepo, likeRepo)
			
			pic, err := service.UploadPicture(context.Background(), userID, tc.image)

			assert.Equal(t, tc.expectedErr, err)
			// Can't do a simple assert.Equal on the picture object because of the random uuid in the url
			if tc.expectedPic != nil {
				assert.NotNil(t, pic)
				assert.Equal(t, tc.expectedPic.UserID, pic.UserID)
				assert.Equal(t, tc.expectedPic.URL, pic.URL)
				assert.Equal(t, tc.expectedPic.IsProfilePic, pic.IsProfilePic)
			} else {
				assert.Nil(t, pic)
			}
		})
	}
}