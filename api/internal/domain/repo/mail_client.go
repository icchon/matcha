package repo

import (
	"context"
)

type MailClient interface {
	SendRawEmail(ctx context.Context, toEmail, subject, htmlBody, textBody string) error
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
