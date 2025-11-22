package server

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/jmoiron/sqlx"

// 	"github.com/icchon/matcha/api/internal/domain/repo"
// 	"github.com/icchon/matcha/api/internal/infrastructure/postgres"
// 	"github.com/icchon/matcha/api/internal/infrastructure/uow"
// 	"github.com/icchon/matcha/api/internal/presentation/handler"
// 	appmiddleware "github.com/icchon/matcha/api/internal/presentation/middleware"
// 	"github.com/icchon/matcha/api/internal/service/auth"
// 	"github.com/icchon/matcha/api/internal/service/mail"
// 	"github.com/icchon/matcha/api/internal/service/user"
// )

// type Config struct {
// 	ServerAddress      string
// 	JWTSigningKey      string
// 	HMACSecretKey      string
// 	GoogleClientID     string
// 	GoogleClientSecret string
// 	GithubClientID     string
// 	GithubClientSecret string
// 	RidirectURI        string

// 	SMTP_HOST     string
// 	SMTP_PORT     string
// 	SMTP_USERNAME string
// 	SMTP_PASSWORD string
// 	SMTP_SENDER   string

// 	BASE_URL string
// }

// type Server struct {
// 	db     *sqlx.DB
// 	router *chi.Mux

// 	config     *Config
// 	httpServer *http.Server
// }

// func NewServer(
// 	db *sqlx.DB,
// 	config *Config,
// 	smtpClient repo.MailClient,
// 	googleClient repo.OAuthClient,
// 	githubClient repo.OAuthClient,
// ) *Server {
// 	unitOfWork := uow.NewUnitOfWork(db)
// 	userRepository := postgres.NewUserRepository(db)
// 	authRepository := postgres.NewAuthRepository(db)
// 	refreshRepository := postgres.NewRefreshTokenRepository(db)
// 	passwordResetRepository := postgres.NewPasswordResetRepository(db)
// 	verificationRepository := postgres.NewVerificationTokenRepository(db)

// 	userService := user.NewUserService(unitOfWork, userRepository)
// 	mailService := mail.NewApplicationMailService(smtpClient, config.BASE_URL)
// 	authService := auth.NewAuthService(unitOfWork, authRepository, userRepository, refreshRepository, passwordResetRepository, verificationRepository, googleClient, githubClient, mailService, config.HMACSecretKey, config.JWTSigningKey)
// 	userHandler := handler.NewUserHandler(userService)
// 	sampleHander := handler.NewSampleHandler()
// 	authHandler := handler.NewAuthHandler(authService)

// 	mux := chi.NewRouter()

// 	server := &Server{
// 		db:     db,
// 		router: mux,
// 		config: config,
// 	}

// 	server.setupRoutes(userHandler, sampleHander, authHandler)

// 	return server
// }

// func (s *Server) setupRoutes(uh *handler.UserHandler, sh *handler.SampleHandler, ah *handler.AuthHandler) {
// 	s.router.Use(middleware.RequestID)
// 	s.router.Use(middleware.Logger)
// 	s.router.Use(middleware.Recoverer)
// 	s.router.Use(middleware.Timeout(60 * time.Second))

// 	s.router.Route("/api/v1", func(r chi.Router) {
// 		r.Get("/sample", sh.GreetingHandler)

// 		r.Route("/auth", func(r chi.Router) {
// 			r.Post("/login", ah.LoginHandler)
// 			r.Post("/signup", ah.SignupHandler)
// 			r.Group(func(r chi.Router) {
// 				r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
// 				r.Post("/logout", ah.LogoutHandler)
// 			})
// 			r.Route("/verify", func(r chi.Router) {
// 				r.Post("/mail", ah.SendVerificationEmailHandler)
// 				r.Get("/{token}", ah.VerifyEmailHandler)
// 			})
// 			r.Route("/oauth", func(r chi.Router) {
// 				r.Route("/google", func(r chi.Router) {
// 					r.Post("/login", ah.GoogleLoginHandler)
// 				})
// 				r.Route("/github", func(r chi.Router) {
// 					r.Post("/login", ah.GithubLoginHandler)
// 				})
// 			})
// 			r.Route("/password", func(r chi.Router) {
// 				r.Post("/forgot", ah.PasswordResetHandler)
// 				r.Post("/reset", ah.PasswordResetConfirmHandler)
// 			})
// 		})

// 	})

// 	log.Println("Routes registered.")
// }

// func (s *Server) Start() error {
// 	s.httpServer = &http.Server{
// 		Addr:         s.config.ServerAddress,
// 		Handler:      s.router,
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 		IdleTimeout:  120 * time.Second,
// 	}

// 	log.Printf("Starting HTTP server on %s", s.config.ServerAddress)
// 	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 		return fmt.Errorf("could not listen on %s: %v", s.config.ServerAddress, err)
// 	}
// 	return nil
// }

// // Shutdown はGraceful Shutdownを行います。
// func (s *Server) Shutdown(ctx context.Context) error {
// 	log.Println("Shutting down server gracefully...")
// 	return s.httpServer.Shutdown(ctx)
// }
