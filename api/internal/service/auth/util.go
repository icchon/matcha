package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

func IsValidPasswordFormat(password string) bool {
	const MIN_PASSWORD_LENGTH int = 8
	return len(password) >= MIN_PASSWORD_LENGTH
}

func IsValidEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateAccessToken(userID uuid.UUID, isVerified bool, authMethod entity.AuthProvider, secretKey string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &entity.AppClaims{
		UserID:     userID,
		IsVerified: isVerified,
		AuthMethod: authMethod,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 'exp'
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 'iat'
			Issuer:    "Matcha",                           // 'iss'
			Subject:   userID.String(),                    // 標準のsubjectも設定
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return tokenString, nil
}

func (s *authService) VerifyRefreshToken(ctx context.Context, tokenString string) (*entity.RefreshToken, error) {
	tokenHash := HashTokenWithHMAC(tokenString, s.hmacSecretKey)
	token, err := s.refreshTokenRepo.Find(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if token == nil {
		return nil, apperrors.ErrUnauthorized
	}
	if token.Revoked {
		return nil, apperrors.ErrUnauthorized
	}
	if token.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrUnauthorized
	}
	return token, nil
}

func GenerateEmailToken() string {
	return uuid.New().String()
}

func GenerateRefreshToken() string {
	return uuid.New().String()
}

func HashTokenWithHMAC(token string, secretKey string) string {
	h := hmac.New(sha512.New, []byte(secretKey))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func CheckTokenWithHMAC(receivedToken, storedHash string, secretKey string) bool {
	expectedHash := HashTokenWithHMAC(receivedToken, secretKey)
	return hmac.Equal([]byte(expectedHash), []byte(storedHash))
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}
