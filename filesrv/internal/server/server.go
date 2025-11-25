package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
)

type ServerConfig struct {
	ServerAddress string
	UploadDir     string
	BaseUrl       string
}

type Server struct {
	router     *chi.Mux
	config     *ServerConfig
	httpServer *http.Server
}

func NewServer(config *ServerConfig) *Server {

	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(config.UploadDir, 0755); err != nil {
			fmt.Printf("Failed to create uploads directory: %v\n", err)
			return nil
		}
	}

	r := chi.NewRouter()

	server := &Server{
		router: r,
		config: config,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(config.UploadDir, config.BaseUrl)
	r.Post("/upload", h.UploadImageHandler)
	r.Get("/dog", h.RandomDogHandler)
	fileServer(r, "/images", http.Dir(config.UploadDir))
	return server
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         s.config.ServerAddress,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", s.config.ServerAddress)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v", s.config.ServerAddress, err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	r.Handle(path+"*", http.StripPrefix(path, http.FileServer(root)))
}
