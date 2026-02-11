package server

import (
	"testing"

	"github.com/icchon/matcha/api/internal/domain/client"
	smtp "github.com/icchon/matcha/api/internal/infrastructure/mail"
)

func TestNewMailClient(t *testing.T) {
	dummyConfig := client.MailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
	}

	tests := []struct {
		name        string
		mode        string
		wantType    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "smtp mode returns SmtpClient",
			mode:     "smtp",
			wantType: "*mail.SmtpClient",
			wantErr:  false,
		},
		{
			name:     "mock mode returns mockMailClient",
			mode:     "mock",
			wantType: "*mail.mockMailClient",
			wantErr:  false,
		},
		{
			name:        "invalid mode returns error",
			mode:        "invalid",
			wantErr:     true,
			errContains: "invalid mail mode",
		},
		{
			name:        "empty mode returns error",
			mode:        "",
			wantErr:     true,
			errContains: "invalid mail mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, err := NewMailClient(tt.mode, dummyConfig)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewMailClient(%q) expected error containing %q, but got nil. Check that invalid modes are rejected.",
						tt.mode, tt.errContains)
				}
				if !containsSubstring(err.Error(), tt.errContains) {
					t.Fatalf("NewMailClient(%q) error = %q, want error containing %q. Check the error message format in NewMailClient.",
						tt.mode, err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewMailClient(%q) unexpected error: %v. Expected a valid %s to be returned.",
					tt.mode, err, tt.wantType)
			}

			switch tt.mode {
			case "smtp":
				if _, ok := mc.(*smtp.SmtpClient); !ok {
					t.Fatalf("NewMailClient(%q) returned type %T, want *mail.SmtpClient. Check that 'smtp' mode creates smtp.NewSmtpClient.",
						tt.mode, mc)
				}
			case "mock":
				// mockMailClient is unexported, so we verify it's not SmtpClient and implements MailClient
				if _, ok := mc.(*smtp.SmtpClient); ok {
					t.Fatalf("NewMailClient(%q) returned %T (*mail.SmtpClient), want mockMailClient. Check that 'mock' mode creates smtp.NewMockMailClient.",
						tt.mode, mc)
				}
				if mc == nil {
					t.Fatalf("NewMailClient(%q) returned nil client. Expected a non-nil mock MailClient.",
						tt.mode)
				}
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
