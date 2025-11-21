package service

import (
	"context"
)

type MailService interface {
	SendVerificationEmail(ctx context.Context, toEmail string, username string, token string) error
	SendPasswordResetEmail(ctx context.Context, toEmail string, token string) error
}
