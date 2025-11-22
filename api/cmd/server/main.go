package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv" // Re-added import
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/icchon/matcha/api/internal/server"
	"github.com/icchon/matcha/api/internal/infrastructure/oauth"
	smtp "github.com/icchon/matcha/api/internal/infrastructure/mail" // Corrected import with alias
	"github.com/icchon/matcha/api/internal/domain/repo"
	)
	
	func checkEnv() error {
		envKeys := []string{
			"SERVER_ADDRESS",
			"DATABASE_URL",
			"JWT_SIGNING_KEY",
			"GOOGLE_CLIENT_ID",
			"GOOGLE_CLIENT_SECRET",
			"GITHUB_CLIENT_ID",
			"GITHUB_CLIENT_SECRET",
			"REDIRECT_URI",
			"HMAC_SECRET_KEY",
			"SMTP_HOST",
			"SMTP_PORT",
			"SMTP_USERNAME",
			"SMTP_PASSWORD",
			"SMTP_SENDER",
			"BASE_URL",
		}
	
		for _, envKey := range envKeys {
			if _, exists := os.LookupEnv(envKey); !exists {
				log.Printf("Missing environment variable: %s", envKey)
				return errors.New("missing environment variable")
			}
		}
		log.Println("All required environment variables are set.")
		return nil
	}
	
	func getEnv(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return ""
	}
	
	func main() {
		if err := checkEnv(); err != nil {
			log.Fatalf("Environment check failed: %v", err)
		}
	
		cfg := &server.Config{
			ServerAddress:      getEnv("SERVER_ADDRESS"),
			JWTSigningKey:      getEnv("JWT_SIGNING_KEY"),
			HMACSecretKey:      getEnv("HMAC_SECRET_KEY"),
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
			GithubClientID:     getEnv("GITHUB_CLIENT_ID"),
			GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET"),
			RidirectURI:        getEnv("REDIRECT_URI"),
			SMTP_HOST:          getEnv("SMTP_HOST"),
			SMTP_PORT:          getEnv("SMTP_PORT"),
			SMTP_USERNAME:      getEnv("SMTP_USERNAME"),
			SMTP_PASSWORD:      getEnv("SMTP_PASSWORD"),
			SMTP_SENDER:        getEnv("SMTP_SENDER"),
			BASE_URL:           getEnv("BASE_URL"),
		}
	
		db, err := sqlx.Connect("postgres", getEnv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
	
		log.Println("Successfully connected to the database.")
	
		// Create concrete client instances
		googleClient := oauth.NewGoogleClient(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.RidirectURI)
		githubClient := oauth.NewGithubClient(cfg.GithubClientID, cfg.GithubClientSecret, cfg.RidirectURI)
	
		smtpPort, err := strconv.Atoi(cfg.SMTP_PORT)
		if err != nil {
			log.Fatalf("Invalid SMTP_PORT: %v", err)
		}
		smtpClient := smtp.NewSmtpClient(repo.MailConfig{
			Host:     cfg.SMTP_HOST,
			Port:     smtpPort,
			Username: cfg.SMTP_USERNAME,
			Password: cfg.SMTP_PASSWORD,
			From:     cfg.SMTP_SENDER,
		})
	
		srv := server.NewServer(db, cfg, smtpClient, googleClient, githubClient)
		if srv == nil {
			log.Fatalf("Faild to setup server")
		}
	go func() {

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received. Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shut down gracefully.")
}
