package main

import (
	"net/http"
	"os"
	"time"
	"log"
	"os/signal"
	"syscall"
	"context"
	"github.com/icchon/matcha/filesrv/internal/server"
)


func main() {
	var config server.ServerConfig
	var ok bool
	config.ServerAddress, ok = os.LookupEnv("SERVER_ADDRESS")
	if !ok{
		log.Fatalf("SERVER_ADDRESS not set")
	}
	config.UploadDir, ok = os.LookupEnv("UPLOAD_DIR")
	if !ok{
		log.Fatalf("UPLOAD_DIR not set")
	}
	config.BaseUrl, ok = os.LookupEnv("BASE_URL")
	if !ok{
		log.Fatalf("BASE_URL not set")
	}
	srv := server.NewServer(&config)
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
