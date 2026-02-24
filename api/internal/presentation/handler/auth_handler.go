package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginHandlerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginHandlerResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	IsVerified   bool      `json:"is_verified"`
	AuthMethod   string    `json:"auth_method"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

// /auth/login POST
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	auth, access, refresh, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, LoginHandlerResponse{AccessToken: access, RefreshToken: refresh, UserID: auth.UserID, IsVerified: auth.IsVerified, AuthMethod: string(auth.Provider)})
}

// auth/logout POST
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	if err := h.authService.Logout(r.Context(), id); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, nil)
}

// auth/signup POST
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.authService.Signup(r.Context(), req.Email, req.Password); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, SignupResponse{Message: "Please check your email to verify your account"})
}

// auth/verify/{token} GET
func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if err := h.authService.VerifyEmail(r.Context(), token); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, nil)
}

// auth/verify POST
type SendVerificationEmailRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}
type SendVerificationEmailResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) SendVerificationEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req SendVerificationEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.authService.SendVerificationEmail(r.Context(), req.Email, req.UserID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, SendVerificationEmailResponse{Message: "Please check your email to verify your account"})
}

// auth/password/forgot POST
type PasswordResetRequest struct {
	Email string `json:"email"`
}
type PasswordResetResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.authService.SendPasswordResetEmail(r.Context(), req.Email); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, PasswordResetResponse{Message: "Please check your email to reset your password"})
}

// auth/password/reset
type PasswordResetConfirmRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *AuthHandler) PasswordResetConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.authService.ConfirmPassword(r.Context(), req.Token, req.Password); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, nil)
}

// auth/oauth/google/login
type GoogleLoginRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}
type GoogleLoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	IsVerified   bool      `json:"is_verified"`
	AuthMethod   string    `json:"auth_method"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func (h *AuthHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	auth, access, refresh, err := h.authService.LoginOAuth(r.Context(), req.Code, req.CodeVerifier, entity.ProviderGoogle)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, GoogleLoginResponse{AccessToken: access, RefreshToken: refresh, UserID: auth.UserID, IsVerified: auth.IsVerified, AuthMethod: string(auth.Provider)})
}

// auth/oauth/github/login
type GithubLoginRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}
type GithubLoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	IsVerified   bool      `json:"is_verified"`
	AuthMethod   string    `json:"auth_method"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func (h *AuthHandler) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	auth, access, refresh, err := h.authService.LoginOAuth(r.Context(), req.Code, req.CodeVerifier, entity.ProviderGithub)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, GoogleLoginResponse{AccessToken: access, RefreshToken: refresh, UserID: auth.UserID, IsVerified: auth.IsVerified, AuthMethod: string(auth.Provider)})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// /auth/refresh POST
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	newAccessToken, err := h.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, RefreshResponse{AccessToken: newAccessToken})
}
