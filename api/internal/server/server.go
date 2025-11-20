package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/icchon/matcha/api/internal/infrastructure/postgres"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
	"github.com/icchon/matcha/api/internal/presentation/handler"
	appmiddleware "github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/icchon/matcha/api/internal/service/user"
)

type Config struct {
	ServerAddress string
	JWTSigningKey string
}

type Server struct {
	db     *sqlx.DB
	router *chi.Mux

	config     *Config
	httpServer *http.Server
}

func NewServer(db *sqlx.DB, config *Config) *Server {
	unitOfWork := uow.NewUnitOfWork(db)
	userRepository := postgres.NewUserRepository(db)
	userService := user.NewUserService(unitOfWork, userRepository)
	userHandler := handler.NewUserHandler(userService)
	sampleHander := handler.NewSampleHandler()

	mux := chi.NewRouter()

	server := &Server{
		db:     db,
		router: mux,
		config: config,
	}

	server.setupRoutes(userHandler, sampleHander)

	return server
}

func (s *Server) setupRoutes(uh *handler.UserHandler, sh *handler.SampleHandler) {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))

				r.Get("/{userID}", uh.FindUserHandler)
			})
		})

		s.router.Get("/sample", sh.GreetingHandler)
	})

	log.Println("Routes registered.")
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

// Shutdown はGraceful Shutdownを行います。
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
