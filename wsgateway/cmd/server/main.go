package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/wsgateway/internal/server"
	"github.com/joho/godotenv"
)

func checkEnv() error {
	envKeys := []string{
		"SERVER_ADDR",
		"REDIS_ADDR",
		"JWT_SIGNING_KEY",
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

	redisAddr := getEnv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}
	defer rdb.Close()
	log.Println("Successfully connected to Redis.")

	conf := &server.ServerConfig{
		ServerAddr:    getEnv("SERVER_ADDR"),
		JwtSigningKey: getEnv("JWT_SIGNING_KEY"),
	}
	srv := server.NewServer(rdb, conf)
	log.Printf("Starting WebSocket server on %s", conf.ServerAddr)

	go func() {
		if err := srv.Start(); err != nil {
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
