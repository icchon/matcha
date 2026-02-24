package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*entity.Auth, string, string, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	Signup(ctx context.Context, email, password string) error
	VerifyEmail(ctx context.Context, token string) error
	SendVerificationEmail(ctx context.Context, email string, userID uuid.UUID) error
	ConfirmPassword(ctx context.Context, token string, password string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	LoginOAuth(ctx context.Context, code string, codeVerifier string, provider entity.AuthProvider) (a *entity.Auth, access string, refresh string, e error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}
