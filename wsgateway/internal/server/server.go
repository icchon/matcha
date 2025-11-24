package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	ServerAddr    string
	JwtSigningKey string
}

type Server struct {
	rdb  *redis.Client
	conf *ServerConfig

	router     *chi.Mux
	httpServer *http.Server
	gateway    *Gateway
}

func NewServer(rdb *redis.Client, conf *ServerConfig) *Server {
	mux := chi.NewRouter()

	gateway := NewGateway(rdb)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(AuthMiddleware(conf.JwtSigningKey))
	mux.HandleFunc("/ws", gateway.handleConnections)
	return &Server{
		rdb:     rdb,
		conf:    conf,
		gateway: gateway,
		router:  mux,
	}
}

func (s *Server) Start() error {
	ctx := context.Background()
	s.gateway.SubscribeChannel(ctx, NotificationChannel, s.gateway.NotificationHandler)
	s.gateway.SubscribeChannel(ctx, ChatOutgoingChannel, s.gateway.ChatMessageHandler)
	s.gateway.SubscribeChannel(ctx, AckChannel, s.gateway.AckHandler)
	s.gateway.SubscribeChannel(ctx, PresenceOutgoingChannel, s.gateway.PresenceHandler)
	s.gateway.SubscribeChannel(ctx, ReadOutgoingChannel, s.gateway.ReadHandler)

	s.httpServer = &http.Server{
		Addr:         s.conf.ServerAddr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
