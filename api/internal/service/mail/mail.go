package mail

import (
	"context"
	"fmt"
	// "github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type applicationMailService struct {
	mailClient client.MailClient
	baseURL    string
}

var _ service.MailService = (*applicationMailService)(nil)

func NewApplicationMailService(client client.MailClient, baseURL string) *applicationMailService {
	return &applicationMailService{
		mailClient: client,
		baseURL:    baseURL,
	}
}

func (s *applicationMailService) SendVerificationEmail(ctx context.Context, toEmail string, username string, token string) error {
	subject := "Matcha: アカウントの確認が必要です"
	verificationLink := fmt.Sprintf("%s/verify?token=%s", s.baseURL, token)

	htmlBody := fmt.Sprintf(`
		<h1>ようこそ、%sさん！</h1>
		<p>以下のリンクをクリックしてアカウントを有効化してください。</p>
		<a href="%s">アカウントを有効化</a>
	`, username, verificationLink)

	return s.mailClient.SendRawEmail(ctx, toEmail, subject, htmlBody, "")
}

func (s *applicationMailService) SendPasswordResetEmail(ctx context.Context, toEmail string, token string) error {
	subject := "Matcha: パスワードのリセット依頼"
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)

	htmlBody := fmt.Sprintf(`
        <h1>パスワードリセット</h1>
        <p>以下のリンクから新しいパスワードを設定してください。</p>
        <a href="%s">パスワードをリセット</a>
    `, resetLink)

	return s.mailClient.SendRawEmail(ctx, toEmail, subject, htmlBody, "")
}
