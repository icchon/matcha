package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type googleClient struct {
	googleClientID     string
	googleClientSecret string
	googleRedirectURL  string
}

type GoogleClaims struct {
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
	Sub string `json:"sub"`

	//email
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`

	//profile
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
}

var _ repo.OAuthClient = (*googleClient)(nil)

func NewGoogleClient(googleClientID string, googleClientSecret string, googleRedirectURL string) *googleClient {
	return &googleClient{
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		googleRedirectURL:  googleRedirectURL,
	}
}

func (c *googleClient) ExchangeCode(ctx context.Context, code string, codeVerifier string) (*repo.OAuthInfo, error) {
	conf := &oauth2.Config{
		ClientID:     c.googleClientID,
		ClientSecret: c.googleClientSecret,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  c.googleRedirectURL,
	}
	ops := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	}

	token, err := conf.Exchange(ctx, code, ops...)
	if err != nil {
		return nil, fmt.Errorf("トークン交換エラー: %w", err)
	}

	idTokenRaw, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("IDトークンがレスポンスに含まれていません")
	}

	payload, err := idtoken.Validate(ctx, idTokenRaw, c.googleClientID)
	if err != nil {
		return nil, fmt.Errorf("IDトークン検証エラー: %w", err)
	}

	var claims GoogleClaims
	claimsJSON, err := json.Marshal(payload.Claims)
	if err != nil {
		return nil, fmt.Errorf("構造体へのマーシャリング失敗: %w", err)
	}

	err = json.Unmarshal(claimsJSON, &claims)
	if err != nil {
		return nil, fmt.Errorf("構造体へのアンマーシャリング失敗: %w", err)
	}

	return &repo.OAuthInfo{
		Sub:           claims.Sub,
		Iss:           claims.Iss,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
	}, nil
}
