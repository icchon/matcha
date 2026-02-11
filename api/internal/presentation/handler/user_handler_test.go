package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_UpdateMyUserDataHandler(t *testing.T) {
	userID := uuid.New()

	updatedData := &entity.UserData{UserID: userID}

	testCases := []struct {
		name           string
		setupMocks     func(mockUserService *mock.MockUserService)
		body           interface{}
		ctx            context.Context
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name: "Success - service is called and re-fetched data returned",
			setupMocks: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().UpdateUserData(gomock.Any(), gomock.Any()).Return(nil)
				mockUserService.EXPECT().GetUserData(gomock.Any(), userID).Return(updatedData, nil)
			},
			body: entity.UserData{},
			ctx:  context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, userID.String(),
					"Response should contain the user ID from re-fetched data, not echoed request body.")
			},
		},
		{
			name: "Service error on update - returns appropriate error response",
			setupMocks: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().UpdateUserData(gomock.Any(), gomock.Any()).Return(apperrors.ErrInternalServer)
			},
			body:           entity.UserData{},
			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Internal server error."}`, body,
					"Service error should propagate. Check that handler calls h.userService.UpdateUserData().")
			},
		},
		{
			name: "Service error on re-fetch - returns appropriate error response",
			setupMocks: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().UpdateUserData(gomock.Any(), gomock.Any()).Return(nil)
				mockUserService.EXPECT().GetUserData(gomock.Any(), userID).Return(nil, apperrors.ErrInternalServer)
			},
			body:           entity.UserData{},
			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Internal server error."}`, body,
					"GetUserData error after successful update should propagate.")
			},
		},
		{
			name:           "Invalid request body - returns 400",
			setupMocks:     func(mockUserService *mock.MockUserService) {},
			body:           `{invalid json`,
			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
			expectedStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Invalid input provided."}`, body,
					"Invalid JSON body should return 400.")
			},
		},
		{
			name:           "No UserID in context - returns 401",
			setupMocks:     func(mockUserService *mock.MockUserService) {},
			body:           entity.UserData{},
			ctx:            context.Background(),
			expectedStatus: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Authentication failed."}`, body,
					"Missing userID in context should return 401.")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mock.NewMockUserService(ctrl)
			mockProfileService := mock.NewMockProfileService(ctrl)
			tc.setupMocks(mockUserService)

			handler := NewUserHandler(mockUserService, mockProfileService)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPut, "/users/me/data", bytes.NewReader(bodyBytes))
			req = req.WithContext(tc.ctx)

			rr := httptest.NewRecorder()
			handler.UpdateMyUserDataHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code,
				"HTTP status mismatch. Check UpdateMyUserDataHandler logic flow.")
			tc.checkBody(t, rr.Body.String())
		})
	}
}
