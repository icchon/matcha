package repo

import (
	"context"
)

type OAuthClient interface {
	ExchangeCode(ctx context.Context, code string, codeVerifier string) (*OAuthInfo, error)
}

type OAuthInfo struct {
	Sub           string `json:"sub"`
	Iss           string `json:"iss"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
