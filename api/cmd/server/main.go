package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/icchon/matcha/api/internal/server"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, loading from environment variables: %v", err)
	}

	if err := checkEnv(); err != nil {
		log.Fatalf("Environment check failed: %v", err)
	}

	cfg := &server.Config{
		ServerAddress: getEnv("SERVER_ADDRESS"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY"),
	}

	db, err := sqlx.Connect("postgres", getEnv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the database.")

	srv := server.NewServer(db, cfg)

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
