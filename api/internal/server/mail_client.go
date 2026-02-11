package server

import (
	"fmt"

	"github.com/icchon/matcha/api/internal/domain/client"
	smtp "github.com/icchon/matcha/api/internal/infrastructure/mail"
)

// NewMailClient creates a MailClient based on the given mode.
// mode must be "smtp" or "mock".
func NewMailClient(mode string, config client.MailConfig) (client.MailClient, error) {
	switch mode {
	case "smtp":
		return smtp.NewSmtpClient(config), nil
	case "mock":
		return smtp.NewMockMailClient(), nil
	default:
		return nil, fmt.Errorf("invalid mail mode: %s (must be 'smtp' or 'mock')", mode)
	}
}
