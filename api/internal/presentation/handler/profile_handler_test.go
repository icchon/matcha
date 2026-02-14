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
	// Minimal valid JPEG magic bytes for http.DetectContentType
	imageData := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46}
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
		{
			name:       "No UserID in context - returns 401",
			setupMocks: func(mockSvc *mock.MockProfileService) {},
			buildRequest: func(t *testing.T) *http.Request {
				body, contentType := createMultipartBody(t, "image", imageData)
				req := httptest.NewRequest(http.MethodPost, "/profile/pictures", body)
				req.Header.Set("Content-Type", contentType)
				return req
			},
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
