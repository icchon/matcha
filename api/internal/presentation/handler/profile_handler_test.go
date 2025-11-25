package handler

// import (
// 	"bytes"
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// 	"github.com/icchon/matcha/api/internal/apperrors"
// 	"github.com/icchon/matcha/api/internal/domain/entity"
// 	"go.uber.org/mock/gomock"

// 	"github.com/icchon/matcha/api/internal/mock"
// 	"github.com/icchon/matcha/api/internal/presentation/middleware"
// 	"github.com/stretchr/testify/assert"
// )

// func TestProfileHandler_CreateProfileHandler(t *testing.T) {
// 	userID := uuid.New()
// 	reqBody := CreateProfileRequest{
// 		Username: sql.NullString{String: "testuser", Valid: true},
// 	}
// 	profile := &entity.UserProfile{
// 		UserID:   userID,
// 		Username: reqBody.Username,
// 	}

// 	testCases := []struct {
// 		name           string
// 		setupMocks     func(mockSvc *mock.MockProfileService)
// 		body           interface{}
// 		ctx            context.Context
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Success",
// 			setupMocks: func(mockSvc *mock.MockProfileService) {
// 				mockSvc.EXPECT().CreateProfile(gomock.Any(), gomock.Any()).Return(profile, nil)
// 			},
// 			body:           reqBody,
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"user_id":"` + userID.String() + `","first_name":{"String":"","Valid":false},"last_name":{"String":"","Valid":false},"username":{"String":"testuser","Valid":true},"gender":{"String":"","Valid":false},"sexual_preference":{"String":"","Valid":false},"biography":{"String":"","Valid":false},"location_name":{"String":"","Valid":false}}`,
// 		},
// 		{
// 			name:           "Invalid JSON",
// 			setupMocks:     func(mockSvc *mock.MockProfileService) {},
// 			body:           `{"bad json`,
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"message":"Invalid input provided."}`,
// 		},
// 		{
// 			name: "Service Error",
// 			setupMocks: func(mockSvc *mock.MockProfileService) {
// 				mockSvc.EXPECT().CreateProfile(gomock.Any(), gomock.Any()).Return(nil, apperrors.ErrInternalServer)
// 			},
// 			body:           reqBody,
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `{"message":"Internal server error."}`,
// 		},
// 		{
// 			name:           "No UserID in Context",
// 			setupMocks:     func(mockSvc *mock.MockProfileService) {},
// 			body:           reqBody,
// 			ctx:            context.Background(),
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `{"message":"Internal server error."}`,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockSvc := mock.NewMockProfileService(ctrl)
// 			tc.setupMocks(mockSvc)

// 			handler := &ProfileHandler{profileSvc: mockSvc}

// 			bodyBytes, _ := json.Marshal(tc.body)
// 			req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(bodyBytes)).WithContext(tc.ctx)
// 			rr := httptest.NewRecorder()

// 			handler.CreateProfileHandler(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code)
// 			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
// 		})
// 	}
// }

// func TestProfileHandler_DeleteProfilePictureHandler(t *testing.T) {
// 	userID := uuid.New()
// 	pictureID := 123

// 	testCases := []struct {
// 		name           string
// 		setupMocks     func(mockSvc *mock.MockProfileService)
// 		ctx            context.Context
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Success",
// 			setupMocks: func(mockSvc *mock.MockProfileService) {
// 				mockSvc.EXPECT().DeletePicture(gomock.Any(), int32(pictureID), userID).Return(nil)
// 			},
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusNoContent,
// 			expectedBody:   `{"message":"Picture deleted successfully"}`,
// 		},
// 		{
// 			name: "Service Error",
// 			setupMocks: func(mockSvc *mock.MockProfileService) {
// 				mockSvc.EXPECT().DeletePicture(gomock.Any(), int32(pictureID), userID).Return(errors.New("some error"))
// 			},
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `{"message":"An unexpected internal error occurred."}`,
// 		},
// 		{
// 			name: "Picture Not Found",
// 			setupMocks: func(mockSvc *mock.MockProfileService) {
// 				mockSvc.EXPECT().DeletePicture(gomock.Any(), int32(pictureID), userID).Return(apperrors.ErrNotFound)
// 			},
// 			ctx:            context.WithValue(context.Background(), middleware.UserIDContextKey, userID),
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody:   `{"message":"The requested resource was not found."}`,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockSvc := mock.NewMockProfileService(ctrl)
// 			tc.setupMocks(mockSvc)

// 			handler := &ProfileHandler{profileSvc: mockSvc}

// 			req := httptest.NewRequest(http.MethodDelete, "/profile/pictures/123", nil)

// 			// Add chi context for URL params
// 			rctx := chi.NewRouteContext()
// 			rctx.URLParams.Add("pictureID", "123")
// 			req = req.WithContext(context.WithValue(tc.ctx, chi.RouteCtxKey, rctx))

// 			rr := httptest.NewRecorder()

// 			handler.DeleteProfilePictureHandler(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code)
// 			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
// 		})
// 	}
// }

// //go:generate mockgen -destination=../../mock/profile_service.go -package=mock github.com/icchon/matcha/api/internal/domain/service ProfileService
