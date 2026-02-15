//go:build integration

package mail_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/infrastructure/mail"
	mailService "github.com/icchon/matcha/api/internal/service/mail"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSendRealEmail(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	mailMode := os.Getenv("MAIL_MODE")
	if mailMode != "smtp" {
		t.Skip("Skipping real email test because MAIL_MODE is not 'smtp'")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpSender := os.Getenv("SMTP_SENDER")
	recipientEmail := os.Getenv("TEST_SMTP_RECIPIENT")

	if smtpHost == "" || smtpPortStr == "" || smtpUsername == "" || smtpPassword == "" || smtpSender == "" {
		t.Fatal("SMTP environment variables are not set")
	}

	if recipientEmail == "" {
		t.Fatal("TEST_SMTP_RECIPIENT environment variable is not set")
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		t.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	mailConfig := client.MailConfig{
		Host:     smtpHost,
		Port:     smtpPort,
		Username: smtpUsername,
		Password: smtpPassword,
		From:     smtpSender,
	}

	smtpClient := mail.NewSmtpClient(mailConfig)
	service := mailService.NewApplicationMailService(smtpClient, "http://localhost:3000")

	err = service.SendVerificationEmail(context.Background(), recipientEmail, "Test User", "test-token-12345")

	assert.NoError(t, err, "Failed to send verification email")
}
