package mail_test

import (
	"context"
	"errors"
	"testing"

	"github.com/icchon/matcha/api/internal/mock"
	mailService "github.com/icchon/matcha/api/internal/service/mail"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSendVerificationEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMailClient := mock.NewMockMailClient(ctrl)
	service := mailService.NewApplicationMailService(mockMailClient, "http://localhost:3100")

	ctx := context.Background()
	toEmail := "test@example.com"
	username := "testuser"
	token := "verification-token-123"

	t.Run("successful verification email sending", func(t *testing.T) {
		mockMailClient.EXPECT().SendRawEmail(
			ctx,
			toEmail,
			gomock.Any(), // Subject
			gomock.Any(), // HTML Body
			gomock.Any(), // Text Body
		).Return(nil).Times(1)

		err := service.SendVerificationEmail(ctx, toEmail, username, token)
		assert.NoError(t, err)
	})

	t.Run("failed verification email sending", func(t *testing.T) {
		mockMailClient.EXPECT().SendRawEmail(
			ctx,
			toEmail,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(errors.New("smtp error")).Times(1)

		err := service.SendVerificationEmail(ctx, toEmail, username, token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "smtp error")
	})
}

func TestSendPasswordResetEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMailClient := mock.NewMockMailClient(ctrl)
	service := mailService.NewApplicationMailService(mockMailClient, "http://localhost:3100")

	ctx := context.Background()
	toEmail := "test@example.com"
	token := "password-reset-token-456"

	t.Run("successful password reset email sending", func(t *testing.T) {
		mockMailClient.EXPECT().SendRawEmail(
			ctx,
			toEmail,
			gomock.Any(), // Subject
			gomock.Any(), // HTML Body
			gomock.Any(), // Text Body
		).Return(nil).Times(1)

		err := service.SendPasswordResetEmail(ctx, toEmail, token)
		assert.NoError(t, err)
	})

	t.Run("failed password reset email sending", func(t *testing.T) {
		mockMailClient.EXPECT().SendRawEmail(
			ctx,
			toEmail,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(errors.New("smtp error")).Times(1)

		err := service.SendPasswordResetEmail(ctx, toEmail, token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "smtp error")
	})
}
