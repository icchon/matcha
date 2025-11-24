package infrastructure

import (
	"context"
	"fmt"
	"github.com/go-mail/mail"
	"github.com/icchon/matcha/api/internal/domain/client"
	"log"
	"time"
)

type SmtpClient struct {
	dialer *mail.Dialer
	from   string
}

var _ client.MailClient = (*SmtpClient)(nil)

func NewSmtpClient(cfg client.MailConfig) *SmtpClient {
	d := mail.NewDialer(
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
	)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return &SmtpClient{
		dialer: d,
		from:   cfg.From,
	}
}

func (c *SmtpClient) SendRawEmail(ctx context.Context, toEmail, subject, htmlBody, textBody string) error {
	m := mail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	// マルチパートメールとして設定
	m.SetBody("text/html", htmlBody)
	if textBody != "" {
		m.AddAlternative("text/plain", textBody)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := c.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("SMTP通信エラー: %w", err)
	}
	log.Println("メール送信成功")
	return nil
}
