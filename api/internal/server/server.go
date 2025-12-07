package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/infrastructure/db/postgres"
	"github.com/icchon/matcha/api/internal/infrastructure/db/uow"
	"github.com/icchon/matcha/api/internal/infrastructure/file"
	smtp "github.com/icchon/matcha/api/internal/infrastructure/mail"
	"github.com/icchon/matcha/api/internal/infrastructure/oauth"
	"github.com/icchon/matcha/api/internal/infrastructure/publisher"
	"github.com/icchon/matcha/api/internal/infrastructure/subscriber"
	"github.com/icchon/matcha/api/internal/presentation/handler"
	appmiddleware "github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/icchon/matcha/api/internal/service/auth"
	"github.com/icchon/matcha/api/internal/service/chat"
	"github.com/icchon/matcha/api/internal/service/mail"
	"github.com/icchon/matcha/api/internal/service/notice"
	"github.com/icchon/matcha/api/internal/service/profile"
	subsvc "github.com/icchon/matcha/api/internal/service/subscriber"
	"github.com/icchon/matcha/api/internal/service/user"
)

type Config struct {
	ServerAddress       string
	JWTSigningKey       string
	HMACSecretKey       string
	GoogleClientID      string
	GoogleClientSecret  string
	GithubClientID      string
	GithubClientSecret  string
	RidirectURI         string
	ImageUploadEndpoint string

	SmtpHost     string
	SmtpPort     string
	SmtpUsername string
	SmtpPassword string
	SmtpSender   string

	BaseUrl string
}

type Server struct {
	router *chi.Mux

	config     *Config
	httpServer *http.Server
}

func NewServer(
	db *sqlx.DB,
	rdb *redis.Client,
	config *Config,
) *Server {
	unitOfWork := uow.NewUnitOfWork(db)

	fileClient := file.NewFilesrvClient(config.ImageUploadEndpoint)
	port, err := strconv.Atoi(config.SmtpPort)
	if err != nil {
		log.Printf("Invalid SMTP_PORT: %v", err)
		return nil
	}
	mailClient := smtp.NewSmtpClient(client.MailConfig{ // Use real SmtpClient
		Host:     config.SmtpHost,
		Port:     port,
		Username: config.SmtpUsername,
		Password: config.SmtpPassword,
		From:     config.SmtpSender,
	})
	githubClient := oauth.NewGithubClient(config.GithubClientID, config.GithubClientSecret, config.RidirectURI)
	googleClient := oauth.NewGoogleClient(config.GoogleClientID, config.GoogleClientSecret, config.RidirectURI)

	notificationPub := publisher.NewNotificationPublisher(rdb)
	ackPub := publisher.NewAckPublisher(rdb)
	presencePub := publisher.NewPresencePublisher(rdb)
	chatPub := publisher.NewChatPublisher(rdb)
	readPub := publisher.NewReadPublisher(rdb)

	userRepository := postgres.NewUserRepository(db)
	authRepository := postgres.NewAuthRepository(db)
	refreshRepository := postgres.NewRefreshTokenRepository(db)
	passwordResetRepository := postgres.NewPasswordResetRepository(db)
	verificationRepository := postgres.NewVerificationTokenRepository(db)
	likeRepository := postgres.NewLikeRepository(db)
	viewRepository := postgres.NewViewRepository(db)
	connectionRepo := postgres.NewConnectionRepository(db)
	profileRepository := postgres.NewUserProfileRepository(db)
	pictureRepository := postgres.NewPictureRepository(db)
	messageRepository := postgres.NewMessageRepository(db)
	notificationRepository := postgres.NewNotificationRepository(db)
	userDataRepository := postgres.NewUserDataRepository(db)
	userTagRepository := postgres.NewUserTagRepository(db)
	tagRepository := postgres.NewTagRepository(db)

	notificationService := notice.NewNotificationService(unitOfWork, notificationRepository, notificationPub)
	userService := user.NewUserService(unitOfWork, likeRepository, viewRepository, connectionRepo, notificationService, userDataRepository, userTagRepository, tagRepository)
	mailService := mail.NewApplicationMailService(mailClient, config.BaseUrl) // Use real mailClient
	authService := auth.NewAuthService(unitOfWork, authRepository, userRepository, refreshRepository, passwordResetRepository, verificationRepository, googleClient, githubClient, mailService, config.HMACSecretKey, config.JWTSigningKey)
	profileService := profile.NewProfileService(unitOfWork, profileRepository, fileClient, pictureRepository, viewRepository, likeRepository, notificationService, userTagRepository, userDataRepository)
	chatService := chat.NewChatService(connectionRepo, messageRepository, profileService)

	userHandler := handler.NewUserHandler(userService, profileService)
	sampleHander := handler.NewSampleHandler()
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	chatHandler := handler.NewChatHandler(chatService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	presenceSub := subscriber.NewPresenceSubscriber(rdb)
	chatSub := subscriber.NewchatSubscriber(rdb)
	readSub := subscriber.NewreadSubscriber(rdb)

	subscHandler := subsvc.NewSubscriberHandler(
		unitOfWork,
		messageRepository,
		readPub,
		ackPub,
		chatPub,
		presencePub,
		userService,
		notificationService,
	)

	subscriberService := subsvc.NewSubscriberService(chatSub, presenceSub, readSub, subscHandler)
	if err := subscriberService.Initialize(context.Background()); err != nil {
		log.Printf("Failed to initialize subscriber service: %v", err)
		return nil
	}

	mux := chi.NewRouter()

	server := &Server{
		router: mux,
		config: config,
	}

	server.setupRoutes(userHandler, sampleHander, authHandler, profileHandler, chatHandler, notificationHandler)

	return server
}

func (s *Server) setupRoutes(uh *handler.UserHandler, sh *handler.SampleHandler, ah *handler.AuthHandler, ph *handler.ProfileHandler, ch *handler.ChatHandler, nh *handler.NotificationHandler) {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/sample", sh.GreetingHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
				r.Post("/logout", ah.LogoutHandler)
			})
			r.Post("/login", ah.LoginHandler)
			r.Post("/signup", ah.SignupHandler)
			r.Post("/verify/mail", ah.SendVerificationEmailHandler)
			r.Get("/verify/{token}", ah.VerifyEmailHandler)
			r.Post("/oauth/google/login", ah.GoogleLoginHandler)
			r.Post("/oauth/github/login", ah.GithubLoginHandler)
			r.Post("/password/forgot", ah.PasswordResetHandler)
			r.Post("/password/reset", ah.PasswordResetConfirmHandler)
			r.Post("/refresh", ah.RefreshAccessTokenHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
				r.Route("/{userID}", func(r chi.Router) {
					r.Post("/like", uh.LikeUserHandler)
					r.Delete("/like", uh.UnlikeUserHandler)
					r.Post("/block", uh.BlockUserHandler)
					r.Delete("/block", uh.UnblockUserHandler)
					r.Get("/profile", ph.GetUserProfileHandler)
					r.Get("/pictures", ph.GetUserPicturesHandler)
				})
			})
		})

		r.Route("/me", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
			r.Delete("/", uh.DeleteMyAccountHandler)
			r.Get("/likes", uh.GetMyLikedListHandler)
			r.Get("/views", uh.GetMyViewedListHandler)
			r.Get("/blocks", uh.GetMyBlockedListHandler)
			r.Get("/chats", ch.GetUserChats)
			r.Get("/notifications", nh.GetUserNotifications)

			r.Route("/data", func(r chi.Router) {
				r.Get("/", uh.GetMyUserDataHandler)
				r.Post("/", uh.CreateMyUserDataHandler)
				r.Put("/", uh.UpdateMyUserDataHandler)
			})

			r.Route("/tags", func(r chi.Router) {
				r.Get("/", uh.GetUserTagsHandler)
				r.Post("/", uh.AddUserTagHandler)
				r.Delete("/{tagID}", uh.DeleteUserTagHandler)
			})

			r.Route("/profile", func(r chi.Router) {
				r.Post("/", ph.CreateProfileHandler)
				r.Put("/", ph.UpdateProfileHandler)
				r.Get("/", ph.GetProfileHandler)
				r.Get("/pictures", ph.GetMyPicturesHandler)
				r.Post("/pictures", ph.UploadProfilePictureHandler)
				r.Delete("/pictures/{pictureID}", ph.DeleteProfilePictureHandler)
				r.Put("/pictures/{pictureID}/status", ph.UpdatePictureStatusHandler)
				r.Get("/likes", ph.GetWhoLikeMeListHandler)
				r.Get("/views", ph.GetWhoViewedMeListHandler)
			})
		})

		r.Route("/tags", func(r chi.Router) {
			r.Get("/", uh.GetAllTagsHandler)
		})
		r.Route("/profiles", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
			r.Get("/", ph.ListProfilesHandler)
			r.Get("/recommends", ph.RecommendProfilesHandler)
		})
		r.Route("/chats/{userID}/messages", func(r chi.Router) {
			r.Use(appmiddleware.AuthMiddleware(s.config.JWTSigningKey))
			r.Get("/", ch.GetChatMessagesHandler)
		})
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