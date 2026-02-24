package notice

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	// "github.com/icchon/matcha/api/internal/domain/client" // Removed: imported and not used
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// MockUOW for notice_test
type mockNoticeUOW struct {
	rm  repo.RepositoryManager
	err error
}

func (m *mockNoticeUOW) Do(ctx context.Context, fn func(rm repo.RepositoryManager) error) error {
	if m.err != nil {
		return m.err
	}
	return fn(m.rm)
}

// MockRM for notice_test
type mockNoticeRM struct {
	repo.RepositoryManager
	notificationRepo *mock.MockNotificationRepository // Use the combined mock
}

func (m *mockNoticeRM) NotificationRepo() repo.NotificationRepository {
	return m.notificationRepo
}

func newUnreadNotification(id int64, recipientID, senderID uuid.UUID) *entity.Notification {
	return &entity.Notification{
		ID:          id,
		RecipientID: recipientID,
		SenderID:    sql.NullString{String: senderID.String(), Valid: true},
		Type:        entity.NotifLike,
		IsRead:      sql.NullBool{Bool: false, Valid: true},
		CreatedAt:   time.Now(),
	}
}

func newReadNotification(id int64, recipientID, senderID uuid.UUID) *entity.Notification {
	return &entity.Notification{
		ID:          id,
		RecipientID: recipientID,
		SenderID:    sql.NullString{String: senderID.String(), Valid: true},
		Type:        entity.NotifLike,
		IsRead:      sql.NullBool{Bool: true, Valid: true},
		CreatedAt:   time.Now(),
	}
}

func TestNotificationService_MarkNotificationAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotificationRepo := mock.NewMockNotificationRepository(ctrl) // Combined mock
	mockPublisher := mock.NewMockPublisher(ctrl)

	uowErr := errors.New("uow error")

	notificationID := int64(1)
	recipientID := uuid.New()
	senderID := uuid.New()

	testCases := []struct {
		name           string
		setupMocks     func()
		uowError       error
		notificationID int64
		recipientID    uuid.UUID
		expectedErr    error
	}{
		{
			name: "Success - Mark unread notification as read",
			setupMocks: func() {
				unReadNotification := newUnreadNotification(notificationID, recipientID, senderID)
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(unReadNotification, nil)
				mockNotificationRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			notificationID: notificationID,
			recipientID:    recipientID,
			expectedErr:    nil,
		},
		{
			name: "Success - Notification already read",
			setupMocks: func() {
				readNotification := newReadNotification(notificationID, recipientID, senderID)
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(readNotification, nil)
			},
			notificationID: notificationID,
			recipientID:    recipientID,
			expectedErr:    nil,
		},
		{
			name: "Notification not found",
			setupMocks: func() {
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(nil, nil)
			},
			notificationID: notificationID,
			recipientID:    recipientID,
			expectedErr:    apperrors.ErrNotFound,
		},
		{
			name: "Recipient ID mismatch",
			setupMocks: func() {
				unReadNotification := newUnreadNotification(notificationID, recipientID, senderID)
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(unReadNotification, nil)
			},
			notificationID: notificationID,
			recipientID:    uuid.New(), // Mismatched recipient
			expectedErr:    apperrors.ErrUnauthorized,
		},
		{
			name: "Find returns DB error",
			setupMocks: func() {
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(nil, errors.New("db error"))
			},
			notificationID: notificationID,
			recipientID:    recipientID,
			expectedErr:    errors.New("db error"),
		},
		{
			name: "UOW fails on update",
			setupMocks: func() {
				unReadNotification := newUnreadNotification(notificationID, recipientID, senderID)
				mockNotificationRepo.EXPECT().Find(gomock.Any(), notificationID).Return(unReadNotification, nil)
				mockNotificationRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(uowErr)
			},
			uowError:       nil, // The UOW itself doesn't fail, the operation within it does
			notificationID: notificationID,
			recipientID:    recipientID,
			expectedErr:    uowErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			mockRM := &mockNoticeRM{notificationRepo: mockNotificationRepo} // Pass the combined mock to RM
			mockUOW := &mockNoticeUOW{rm: mockRM, err: tc.uowError}

			service := NewNotificationService(mockUOW, mockNotificationRepo, mockPublisher) // Pass combined mock as query repo
			err := service.MarkNotificationAsRead(context.Background(), tc.notificationID, tc.recipientID)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
