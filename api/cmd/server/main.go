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

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/icchon/matcha/api/internal/server"
)

func checkEnv() error {
	envKeys := []string{
		"SERVER_ADDR",
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
		"REDIS_ADDR",
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
		log.Println("No .env file found or error loading .env file, proceeding with system environment variables.")
	}
	if err := checkEnv(); err != nil {
		log.Fatalf("Environment check failed: %v", err)
	}

	cfg := &server.Config{
		ServerAddress:       getEnv("SERVER_ADDR"),
		JWTSigningKey:       getEnv("JWT_SIGNING_KEY"),
		HMACSecretKey:       getEnv("HMAC_SECRET_KEY"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET"),
		GithubClientID:      getEnv("GITHUB_CLIENT_ID"),
		GithubClientSecret:  getEnv("GITHUB_CLIENT_SECRET"),
		RidirectURI:         getEnv("REDIRECT_URI"),
		SmtpHost:            getEnv("SMTP_HOST"),
		SmtpPort:            getEnv("SMTP_PORT"),
		SmtpUsername:        getEnv("SMTP_USERNAME"),
		SmtpPassword:        getEnv("SMTP_PASSWORD"),
		SmtpSender:          getEnv("SMTP_SENDER"),
		BaseUrl:             getEnv("BASE_URL"),
		ImageUploadEndpoint: getEnv("IMAGE_UPLOAD_ENDPOINT"),
	}

	db, err := sqlx.Connect("postgres", getEnv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database.")

	redisAddr := getEnv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}
	defer rdb.Close()
	log.Println("Successfully connected to Redis.")

	srv := server.NewServer(db, cfg)
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
