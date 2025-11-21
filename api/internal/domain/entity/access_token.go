package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppClaims struct {
	UserID               uuid.UUID    `json:"sub"` // Standard 'sub' claim for subject
	IsVerified           bool         `json:"is_verified"`
	AuthMethod           AuthProvider `json:"auth_method"`
	jwt.RegisteredClaims              // JWTの標準クレーム (iss, exp, iatなど) を継承
}
