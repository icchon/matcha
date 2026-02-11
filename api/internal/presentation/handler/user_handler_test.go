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

	testCases := []struct {
		name           string
		setupMocks     func(mockUserService *mock.MockUserService)
		body           interface{}
		ctx            context.Context
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name: "Success - service is called and updated data returned",
			setupMocks: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().UpdateUserData(gomock.Any(), gomock.Any()).Return(nil)
			},
			body: entity.UserData{},
			ctx:  context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, userID.String(),
					"Response should contain the user ID. Check that UpdateMyUserDataHandler sets userData.UserID before responding.")
			},
		},
		{
			name: "Service error - returns appropriate error response",
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

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// 	"github.com/icchon/matcha/api/internal/apperrors"
// 	"github.com/icchon/matcha/api/internal/domain/entity"
// 	"github.com/icchon/matcha/api/internal/mock"
// 	"github.com/icchon/matcha/api/internal/presentation/middleware"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// func TestUserHandler_LikeUserHandler(t *testing.T) {
// 	likerID := uuid.New()
// 	likedID := uuid.New()

// 	testCases := []struct {
// 		name           string
// 		setupMocks     func(mockUserService *mock.MockUserService)
// 		likedID        string
// 		ctx            context.Context
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Success - Like",
// 			setupMocks: func(mockUserService *mock.MockUserService) {
// 				mockUserService.EXPECT().LikeUser(gomock.Any(), likerID, likedID).Return(nil, nil)
// 			},
// 			likedID:        likedID.String(),
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"connection":null, "message":"User liked successfully"}`,
// 		},
// 		{
// 			name: "Success - Match",
// 			setupMocks: func(mockUserService *mock.MockUserService) {
// 				conn := &entity.Connection{User1ID: likerID, User2ID: likedID}
// 				mockUserService.EXPECT().LikeUser(gomock.Any(), likerID, likedID).Return(conn, nil)
// 			},
// 			likedID:        likedID.String(),
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"connection":{"user1_id":"` + likerID.String() + `","user2_id":"` + likedID.String() + `", "created_at":"0001-01-01T00:00:00Z"}, "message":"It's a match!"}`,
// 		},
// 		{
// 			name:           "Invalid LikedID",
// 			setupMocks:     func(mockUserService *mock.MockUserService) {},
// 			likedID:        "invalid-uuid",
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"message":"Invalid input provided."}`,
// 		},
// 		{
// 			name: "Service Error",
// 			setupMocks: func(mockUserService *mock.MockUserService) {
// 				mockUserService.EXPECT().LikeUser(gomock.Any(), likerID, likedID).Return(nil, apperrors.ErrInternalServer)
// 			},
// 			likedID:        likedID.String(),
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `{"message":"Internal server error."}`,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockUserService := mock.NewMockUserService(ctrl)
// 			mockProfileService := mock.NewMockProfileService(ctrl)
// 			tc.setupMocks(mockUserService)

// 			handler := NewUserHandler(mockUserService, mockProfileService)

// 			req := httptest.NewRequest(http.MethodPost, "/users/"+tc.likedID+"/like", nil)

// 			rctx := chi.NewRouteContext()
// 			rctx.URLParams.Add("userID", tc.likedID)
// 			req = req.WithContext(context.WithValue(tc.ctx, chi.RouteCtxKey, rctx))

// 			rr := httptest.NewRecorder()
// 			handler.LikeUserHandler(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code)
// 			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
// 		})
// 	}
// }

// func TestUserHandler_UnlikeUserHandler(t *testing.T) {
// 	likerID := uuid.New()
// 	likedID := uuid.New()

// 	testCases := []struct {
// 		name           string
// 		setupMocks     func(mockUserService *mock.MockUserService)
// 		likedID        string
// 		ctx            context.Context
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Success",
// 			setupMocks: func(mockUserService *mock.MockUserService) {
// 				mockUserService.EXPECT().UnlikeUser(gomock.Any(), likerID, likedID).Return(nil)
// 			},
// 			likedID:        likedID.String(),
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"message":"User unliked successfully"}`,
// 		},
// 		{
// 			name:           "Invalid LikedID",
// 			setupMocks:     func(mockUserService *mock.MockUserService) {},
// 			likedID:        "invalid-uuid",
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"message":"Invalid input provided."}`,
// 		},
// 		{
// 			name: "Service Error",
// 			setupMocks: func(mockUserService *mock.MockUserService) {
// 				mockUserService.EXPECT().UnlikeUser(gomock.Any(), likerID, likedID).Return(apperrors.ErrNotFound)
// 			},
// 			likedID:        likedID.String(),
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, likerID),
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody:   `{"message":"The requested resource was not found."}`,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockUserService := mock.NewMockUserService(ctrl)
// 			mockProfileService := mock.NewMockProfileService(ctrl)
// 			tc.setupMocks(mockUserService)

// 			handler := NewUserHandler(mockUserService, mockProfileService)

// 			req := httptest.NewRequest(http.MethodDelete, "/users/"+tc.likedID+"/like", nil)

// 			rctx := chi.NewRouteContext()
// 			rctx.URLParams.Add("userID", tc.likedID)
// 			req = req.WithContext(context.WithValue(tc.ctx, chi.RouteCtxKey, rctx))

// 			rr := httptest.NewRecorder()
// 			handler.UnlikeUserHandler(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code)
// 			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
// 		})
// 	}
// }

// //go:generate mockgen -destination=../../mock/user_service.go -package=mock github.com/icchon/matcha/api/internal/domain/service UserService
