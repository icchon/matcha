package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// 	"go.uber.org/mock/gomock"

// 	"github.com/icchon/matcha/api/internal/mock"
// 	"github.com/icchon/matcha/api/internal/server"
// 	"github.com/icchon/matcha/api/internal/domain/repo"
// )

// // TestMain is used to set up and tear down a test server for integration/API tests.
// // This function can be used by external tools like "api dog" to start the server
// // with mocked external dependencies.
// func TestMain(m *testing.M) {
// 	// Load environment variables (optional, for test config if needed)
// 	// --- Setup for Test Server with Mocks ---
// 	// Create a gomock controller for managing mocks
// 	ctrl := gomock.NewController(nil) // nil for production, or a *testing.T in a normal test func
// 	defer ctrl.Finish()

// 	// Create mock clients
// 	mockMailClient := mock.NewMockMailClient(ctrl)
// 	mockGoogleClient := mock.NewMockOAuthClient(ctrl)
// 	mockGithubClient := mock.NewMockOAuthClient(ctrl)

// 	// Configure mock behavior for integration tests
// 	// For example, make SendRawEmail always return nil (no error)
// 	mockMailClient.EXPECT().SendRawEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
// 	// For OAuth clients, you might expect ExchangeCode and return a dummy OAuthInfo
// 	mockGoogleClient.EXPECT().ExchangeCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(&repo.OAuthInfo{
// 		Sub:           "mock-google-sub",
// 		Email:         "mock-google@example.com",
// 		EmailVerified: true,
// 	}, nil).AnyTimes()
// 	mockGithubClient.EXPECT().ExchangeCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(&repo.OAuthInfo{
// 		Sub:           "mock-github-sub",
// 		Email:         "mock-github@example.com",
// 		EmailVerified: true,
// 	}, nil).AnyTimes()

// 	// Use dummy config values for the test server
// 	cfg := &server.Config{
// 		ServerAddress:      ":3100", // Use a different port for testing, e.g., ":8081"
// 		JWTSigningKey:      "test-jwt-signing-key",
// 		HMACSecretKey:      "test-hmac-secret-key",
// 		GoogleClientID:     "test-google-client-id",
// 		GoogleClientSecret: "test-google-client-secret",
// 		GithubClientID:     "test-github-client-id",
// 		GithubClientSecret: "test-github-client-secret",
// 		RidirectURI:        "http://localhost:8081/auth/oauth/callback", // Test redirect URI
// 		SMTP_HOST:          "smtp.test.com",
// 		SMTP_PORT:          "5252",
// 		SMTP_USERNAME:      "testuser",
// 		SMTP_PASSWORD:      "testpass",
// 		SMTP_SENDER:        "test@test.com",
// 		BASE_URL:           "http://localhost:8081",
// 	}

// 	// For a test, you would typically provide a mock DB connection here or an in-memory DB.
// 	// For this example, we'll use a dummy DB connection for compilation.
// 	db, err := sqlx.Connect("postgres", "postgresql://user:p@localhost:5432/matcha?sslmode=disable") // Use a test DB URL
// 	if err != nil {
// 		log.Fatalf("Failed to connect to test database: %v", err)
// 	}
// 	defer db.Close()

// 	log.Println("Successfully connected to the test database.")

// 	// Instantiate the server with mocked clients
// 	testSrv := server.NewServer(db, cfg, mockMailClient, mockGoogleClient, mockGithubClient)
// 	if testSrv == nil {
// 		log.Fatalf("Failed to setup test server")
// 	}

// 	// --- Start the Test Server ---
// 	go func() {
// 		log.Printf("Starting test HTTP server on %s", cfg.ServerAddress)
// 		if err := testSrv.Start(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Test server failed to start: %v", err)
// 		}
// 	}()

// 	// Wait for server to start (optional, but good practice for integration tests)
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Println("Shutdown signal received. Starting graceful shutdown...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := testSrv.Shutdown(ctx); err != nil {
// 		log.Fatalf("Server forced to shutdown: %v", err)
// 	}
// 	log.Println("Test server shut down gracefully.")
// 	os.Exit(0)
// }
