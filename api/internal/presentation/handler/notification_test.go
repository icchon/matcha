package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationHandler_MarkNotificationAsReadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifService := mock.NewMockNotificationService(ctrl)
	handler := NewNotificationHandler(mockNotifService)

	notificationID := int64(123)
	recipientID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockNotifService.EXPECT().MarkNotificationAsRead(gomock.Any(), notificationID, recipientID).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/me/notifications/123/read", nil)
		ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("Invalid Notification ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/me/notifications/invalid/read", nil)
		ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedBody := `{"message":"Invalid input provided."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Unauthorized - No UserID in Context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/me/notifications/123/read", nil)
		// No UserIDContextKey in context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		expectedBody := `{"message":"Authentication failed."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Service returns Unauthorized Error", func(t *testing.T) {
		mockNotifService.EXPECT().MarkNotificationAsRead(gomock.Any(), notificationID, recipientID).Return(apperrors.ErrUnauthorized)

		req := httptest.NewRequest(http.MethodPut, "/me/notifications/123/read", nil)
		ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		expectedBody := `{"message":"Authentication failed."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Service returns NotFound Error", func(t *testing.T) {
		mockNotifService.EXPECT().MarkNotificationAsRead(gomock.Any(), notificationID, recipientID).Return(apperrors.ErrNotFound)

		req := httptest.NewRequest(http.MethodPut, "/me/notifications/123/read", nil)
		ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		expectedBody := `{"message":"The requested resource was not found."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Service returns InternalServer Error", func(t *testing.T) {
		mockNotifService.EXPECT().MarkNotificationAsRead(gomock.Any(), notificationID, recipientID).Return(errors.New("db error")) // generic error

		req := httptest.NewRequest(http.MethodPut, "/me/notifications/123/read", nil)
		ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()
		handler.MarkNotificationAsReadHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedBody := `{"message":"An unexpected internal error occurred."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})
}

func TestNotificationHandler_MarkAllNotificationsAsReadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifService := mock.NewMockNotificationService(ctrl)
	handler := NewNotificationHandler(mockNotifService)

	recipientID := uuid.New()

	testCases := []struct {
		name           string
		setupMocks     func()
		ctx            context.Context
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMocks: func() {
				mockNotifService.EXPECT().MarkAllNotificationsAsRead(gomock.Any(), recipientID).Return(nil)
			},
			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID),
			expectedStatus: http.StatusNoContent,
			expectedBody:   ``, // No content for 204
		},
		{
			name:           "Unauthorized - No UserID in Context",
			setupMocks:     func() {},
			ctx:            context.Background(),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Authentication failed."}`,
		},
		{
			name: "Service returns InternalServer Error",
			setupMocks: func() {
				mockNotifService.EXPECT().MarkAllNotificationsAsRead(gomock.Any(), recipientID).Return(errors.New("db error"))
			},
			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, recipientID),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"An unexpected internal error occurred."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/me/notifications/read", nil)
			req = req.WithContext(tc.ctx)

			rr := httptest.NewRecorder()
			handler.MarkAllNotificationsAsReadHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, rr.Body.String())
		})
	}
}
