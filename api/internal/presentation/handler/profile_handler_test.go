package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProfileHandler_UploadProfilePictureHandler(t *testing.T) {
	userID := uuid.New()
	imageData := []byte("fake-image-data")
	expectedPicture := &entity.Picture{
		ID:        1,
		UserID:    userID,
		URL:       "https://example.com/pic.jpg",
		CreatedAt: time.Now(),
	}

	createMultipartBody := func(t *testing.T, fieldName string, data []byte) (*bytes.Buffer, string) {
		t.Helper()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(fieldName, "test.jpg")
		assert.NoError(t, err, "Failed to create multipart form file.")
		_, err = io.Copy(part, bytes.NewReader(data))
		assert.NoError(t, err, "Failed to write to multipart form file.")
		writer.Close()
		return body, writer.FormDataContentType()
	}

	testCases := []struct {
		name           string
		setupMocks     func(mockSvc *mock.MockProfileService)
		buildRequest   func(t *testing.T) *http.Request
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name: "Success - upload picture",
			setupMocks: func(mockSvc *mock.MockProfileService) {
				mockSvc.EXPECT().UploadPicture(gomock.Any(), userID, imageData).Return(expectedPicture, nil)
			},
			buildRequest: func(t *testing.T) *http.Request {
				body, contentType := createMultipartBody(t, "image", imageData)
				req := httptest.NewRequest(http.MethodPost, "/profile/pictures", body)
				req.Header.Set("Content-Type", contentType)
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserIDContextKey, userID))
				return req
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, expectedPicture.URL,
					"Response should contain the picture URL. Check that UploadProfilePictureHandler returns the service result.")
			},
		},
		{
			name:       "No file in request - returns error",
			setupMocks: func(mockSvc *mock.MockProfileService) {},
			buildRequest: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/profile/pictures", nil)
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserIDContextKey, userID))
				return req
			},
			expectedStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "message",
					"Missing file should return an error response with a message field.")
			},
		},
		{
			name: "Service error - returns 500",
			setupMocks: func(mockSvc *mock.MockProfileService) {
				mockSvc.EXPECT().UploadPicture(gomock.Any(), userID, imageData).Return(nil, apperrors.ErrInternalServer)
			},
			buildRequest: func(t *testing.T) *http.Request {
				body, contentType := createMultipartBody(t, "image", imageData)
				req := httptest.NewRequest(http.MethodPost, "/profile/pictures", body)
				req.Header.Set("Content-Type", contentType)
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserIDContextKey, userID))
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Internal server error."}`, body,
					"Service error should propagate to response.")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mock.NewMockProfileService(ctrl)
			tc.setupMocks(mockSvc)

			handler := &ProfileHandler{profileSvc: mockSvc}

			req := tc.buildRequest(t)
			rr := httptest.NewRecorder()

			handler.UploadProfilePictureHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code,
				"HTTP status mismatch. Check UploadProfilePictureHandler implementation.")
			tc.checkBody(t, rr.Body.String())
		})
	}
}

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
