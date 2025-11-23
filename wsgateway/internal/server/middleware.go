package server

import (
	"context"
	"net/http"
	"strings"

	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"time"
)

// jwt認証用のendpointをAPIに追加しても良いが面倒くさい上にgateway -> api の依存があるのがなんか嫌なのでこちらで認証

type AppClaims struct {
	UserID               uuid.UUID `json:"sub"` // Standard 'sub' claim for subject
	IsVerified           bool      `json:"is_verified"`
	AuthMethod           string    `json:"auth_method"`
	jwt.RegisteredClaims           // JWTの標準クレーム (iss, exp, iatなど) を継承
}

type ContextKey string

const (
	UserIDContextKey     ContextKey = "userID"
	IsVerifiedContextKey ContextKey = "isVerified"
	AuthMethodContextKey ContextKey = "authMethod"
)

func AuthMiddleware(jwtSigningKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			claims, err := VerifyAccessToken(tokenString, jwtSigningKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, IsVerifiedContextKey, claims.IsVerified)
			ctx = context.WithValue(ctx, AuthMethodContextKey, claims.AuthMethod)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func VerifyAccessToken(tokenString string, secretKey string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("token parsing error: %v", err)
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	if !token.Valid {
		log.Printf("token is invalid")
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*AppClaims)
	if !ok {
		log.Printf("invalid token claims format")
		return nil, fmt.Errorf("invalid token claims format")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		log.Printf("token has expired")
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
