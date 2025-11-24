package mail

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/client"
)

var _ client.MailClient = (*mockMailClient)(nil)

type mockMailClient struct {
}

func NewMockMailClient() *mockMailClient {
	return &mockMailClient{}
}

func (c *mockMailClient) SendRawEmail(ctx context.Context, toEmail, subject, htmlBody, textBody string) error {
	return nil
}
