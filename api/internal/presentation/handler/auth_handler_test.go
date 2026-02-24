package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_LoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	userID := uuid.New()
	auth := &entity.Auth{
		UserID:     userID,
		IsVerified: true,
		Provider:   entity.ProviderLocal,
	}
	accessToken := "new-access-token"
	refreshToken := "new-refresh-token"

	t.Run("Success", func(t *testing.T) {
		mockAuthService.EXPECT().Login(gomock.Any(), "test@example.com", "password").Return(auth, accessToken, refreshToken, nil)

		reqBody := LoginHandlerRequest{Email: "test@example.com", Password: "password"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var res LoginHandlerResponse
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, accessToken, res.AccessToken)
		assert.Equal(t, refreshToken, res.RefreshToken)
	})
}

func TestAuthHandler_RefreshHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	refreshToken := "some-refresh-token"
	newAccessToken := "new-access-token"

	t.Run("Success", func(t *testing.T) {
		mockAuthService.EXPECT().RefreshAccessToken(gomock.Any(), refreshToken).Return(newAccessToken, nil)

		reqBody := RefreshRequest{RefreshToken: refreshToken}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.RefreshHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expectedBody := `{"access_token":"` + newAccessToken + `"}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Service returns unauthorized error", func(t *testing.T) {
		mockAuthService.EXPECT().RefreshAccessToken(gomock.Any(), refreshToken).Return("", apperrors.ErrUnauthorized)

		reqBody := RefreshRequest{RefreshToken: refreshToken}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.RefreshHandler(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		expectedBody := `{"message":"Authentication failed."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader([]byte(`{"bad json`)))
		rr := httptest.NewRecorder()

		handler.RefreshHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedBody := `{"message":"Invalid input provided."}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})
}
