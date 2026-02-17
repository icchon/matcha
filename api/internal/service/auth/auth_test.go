package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// MockUOW for auth_test
type mockAuthUOW struct {
	rm  repo.RepositoryManager
	err error
}

func (m *mockAuthUOW) Do(ctx context.Context, fn func(rm repo.RepositoryManager) error) error {
	if m.err != nil {
		return m.err
	}
	return fn(m.rm)
}

// MockRM for auth_test
type mockAuthRM struct {
	repo.RepositoryManager
	refreshTokenRepo repo.RefreshTokenRepository
	userRepo         repo.UserRepository
}

func (m *mockAuthRM) RefreshTokenRepo() repo.RefreshTokenRepository {
	return m.refreshTokenRepo
}

func (m *mockAuthRM) UserRepo() repo.UserRepository {
	return m.userRepo
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock.NewMockAuthQueryRepository(ctrl)
	mockUserQueryRepo := mock.NewMockUserQueryRepository(ctrl)
	mockRefreshTokenQueryRepo := mock.NewMockRefreshTokenQueryRepository(ctrl)
	mockPasswordResetQueryRepo := mock.NewMockPasswordResetQueryRepository(ctrl)
	mockVerificationTokenQueryRepo := mock.NewMockVerificationTokenQueryRepository(ctrl)
	mockMailService := mock.NewMockMailService(ctrl) // Use service.MailService mock
	mockGoogleClient := mock.NewMockOAuthClient(ctrl)
	mockGithubClient := mock.NewMockOAuthClient(ctrl)

	// Unified Repo Mocks (implementing both Query and Command interfaces)
	mockRefreshTokenRepo := mock.NewMockRefreshTokenRepository(ctrl) // Use unified mock
	mockUserRepo := mock.NewMockUserRepository(ctrl)                 // Use unified mock

	userID := uuid.New()
	tokenHash := "some_token_hash"
	expiresAt := time.Now().Add(time.Hour)

	activeRefreshToken := &entity.RefreshToken{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	user := &entity.User{
		ID:        userID,
		CreatedAt: time.Now(),
		LastConnection: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	testCases := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Success - Refresh token revoked and user last connection updated",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.RefreshToken{activeRefreshToken}, nil)
				mockUserQueryRepo.EXPECT().Find(gomock.Any(), userID).Return(user, nil)
				mockRefreshTokenRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "RefreshTokenRepo.Query returns error",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectedErr: apperrors.ErrInternalServer,
		},
		{
			name: "RefreshToken not found",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.RefreshToken{}, nil)
			},
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "UserRepo.Find returns error",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.RefreshToken{activeRefreshToken}, nil)
				mockUserQueryRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, errors.New("db error"))
			},
			expectedErr: apperrors.ErrInternalServer,
		},
		{
			name: "User not found",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.RefreshToken{activeRefreshToken}, nil)
				mockUserQueryRepo.EXPECT().Find(gomock.Any(), userID).Return(nil, nil)
			},
			expectedErr: apperrors.ErrNotFound,
		},
		{
			name: "UOW update RefreshToken fails",
			setupMocks: func() {
				mockRefreshTokenQueryRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.RefreshToken{activeRefreshToken}, nil)
				mockUserQueryRepo.EXPECT().Find(gomock.Any(), userID).Return(user, nil)
				mockRefreshTokenRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("uow update error"))
			},
			expectedErr: errors.New("uow update error"), // Expect the UOW error to be returned
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			mockRM := &mockAuthRM{
				refreshTokenRepo: mockRefreshTokenRepo,
				userRepo:         mockUserRepo,
			}
			mockUOW := &mockAuthUOW{rm: mockRM}

			service := NewAuthService(
				mockUOW,
				mockAuthRepo,
				mockUserQueryRepo,
				mockRefreshTokenQueryRepo,
				mockPasswordResetQueryRepo,
				mockVerificationTokenQueryRepo,
				mockGoogleClient,
				mockGithubClient,
				mockMailService,
				"dummy_hmac_key",
				"dummy_jwt_key",
			)

			err := service.Logout(context.Background(), userID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
